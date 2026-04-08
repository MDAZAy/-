param(
    [string]$BaseUrl = "http://127.0.0.1:8080"
)

$ErrorActionPreference = "Stop"

$health = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
$user = Invoke-RestMethod -Uri "$BaseUrl/api/v1/users/ensure" -Method Post -ContentType "application/json" -Body '{"telegram_id":123456789,"username":"smoke_user","full_name":"Smoke Test"}'
$plans = Invoke-RestMethod -Uri "$BaseUrl/api/v1/plans" -Method Get

if (-not $plans -or $plans.Count -lt 1) {
    throw "Smoke failed: no plans returned."
}

$paymentBody = @{
    user_id = $user.id
    plan_id = $plans[0].id
} | ConvertTo-Json

$payment = Invoke-RestMethod -Uri "$BaseUrl/api/v1/payments/create" -Method Post -ContentType "application/json" -Body $paymentBody
$null = Invoke-WebRequest -Uri "$BaseUrl/mock/payments/$($payment.external_payment_id)" -Method Get
$paymentPage = Invoke-WebRequest -Uri "$BaseUrl/mock/payments/$($payment.external_payment_id)/succeed" -Method Post
$subscription = Invoke-RestMethod -Uri "$BaseUrl/api/v1/subscriptions/active/$($user.id)" -Method Get
$vpnBody = @{ user_id = $user.id } | ConvertTo-Json
$vpn = Invoke-RestMethod -Uri "$BaseUrl/api/v1/vpn/issue" -Method Post -ContentType "application/json" -Body $vpnBody

[PSCustomObject]@{
    health_status         = $health.status
    user_id               = $user.id
    plans_count           = $plans.Count
    payment_status        = $payment.status
    payment_page_status   = $paymentPage.StatusCode
    subscription_status   = $subscription.status
    subscription_plan_id  = $subscription.plan_id
    vpn_provider          = $vpn.provider
    vpn_url               = $vpn.access_url
} | ConvertTo-Json -Depth 4
