// static/script.js

const ws = new WebSocket("ws://localhost:8080/ws");

// When connected
ws.onopen = () => {
  console.log("Connected to Go WebSocket 🟢");
  const status = document.getElementById("status");
  if (status) status.textContent = "🟢 Live connection active";
};

// Handle incoming stock data
ws.onmessage = (event) => {
  const stocks = JSON.parse(event.data);
  updateDashboard(stocks);
};

// On error
ws.onerror = (err) => {
  console.error("WebSocket error:", err);
  const status = document.getElementById("status");
  if (status) status.textContent = "🔴 Connection error";
};

// Update the dashboard
function updateDashboard(stocks) {
  const container = document.getElementById("stocks");
  container.innerHTML = ""; // clear old data

  stocks.forEach((stock) => {
    const card = document.createElement("div");
    card.className = "stock-card";

    const isPositive = stock.change >= 0;
    const color = isPositive ? "#16a34a" : "#dc2626"; // green/red
    const sign = isPositive ? "+" : "";

    card.innerHTML = `
      <div class="symbol">${stock.symbol}</div>
      <div class="price">$${stock.price.toFixed(2)}</div>
      <div class="change" style="color:${color}">
        ${sign}${stock.change.toFixed(2)}%
      </div>
      <div class="time">${stock.time}</div>
    `;

    container.appendChild(card);
  });
}
