const promptEl = document.getElementById("prompt");
const resultEl = document.getElementById("result");
const healthEl = document.getElementById("health-result");

async function requestJSON(url, options = {}) {
  const response = await fetch(url, options);
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || `Request failed: ${response.status}`);
  }
  return data;
}

document.getElementById("generate").addEventListener("click", async () => {
  resultEl.textContent = "Генерация...";
  try {
    const data = await requestJSON("/generate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ prompt: promptEl.value }),
    });
    resultEl.textContent = JSON.stringify(data, null, 2);
  } catch (error) {
    resultEl.textContent = error.message;
  }
});

document.getElementById("health").addEventListener("click", async () => {
  healthEl.textContent = "Проверка...";
  try {
    const data = await requestJSON("/health");
    healthEl.textContent = JSON.stringify(data, null, 2);
  } catch (error) {
    healthEl.textContent = error.message;
  }
});
