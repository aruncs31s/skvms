const API_BASE = "/api";
let deviceId = null;
let voltageChart = null;
let autoRefreshInterval = null;

// Initialize on page load
document.addEventListener("DOMContentLoaded", () => {
  loadNavbar();
  checkAuth();

  // Get device ID from URL
  const path = window.location.pathname;
  const match = path.match(/\/devices\/(\d+)/);
  if (match) {
    deviceId = match[1];
    loadDeviceData();
    loadReadings();
    loadStateHistory();
    
    // Auto-refresh every 30 seconds
    autoRefreshInterval = setInterval(() => {
      loadDeviceData();
      loadReadings();
    }, 30000);
  }

  setupEventListeners();
});

function setupEventListeners() {
  // Control buttons
  document.getElementById("btnTurnOn")?.addEventListener("click", () => controlDevice(4));
  document.getElementById("btnTurnOff")?.addEventListener("click", () => controlDevice(5));
  document.getElementById("btnConfigure")?.addEventListener("click", () => controlDevice(6));
  document.getElementById("btnRefresh")?.addEventListener("click", () => {
    loadDeviceData();
    loadReadings();
    loadStateHistory();
  });

  // Refresh buttons
  document.getElementById("refreshReadings")?.addEventListener("click", loadReadings);
  document.getElementById("refreshHistory")?.addEventListener("click", loadStateHistory);

  // Export CSV
  document.getElementById("exportCsv")?.addEventListener("click", exportToCSV);

  // Readings limit change
  document.getElementById("readingsLimit")?.addEventListener("change", loadReadings);
}

async function loadDeviceData() {
  try {
    const response = await fetch(`${API_BASE}/devices/${deviceId}`);
    if (!response.ok) throw new Error("Failed to load device");

    const data = await response.json();
    const device = data.device;

    // Update title and meta
    document.getElementById("deviceTitle").textContent = device.name || "Device Dashboard";
    document.getElementById("deviceMeta").textContent = `${device.type} • ${device.city}`;

    // Update device info
    document.getElementById("deviceType").textContent = device.type;
    document.getElementById("deviceIP").textContent = device.ip_address;
    document.getElementById("deviceMAC").textContent = device.mac_address;
    document.getElementById("deviceFirmware").textContent = device.firmware_version;
    document.getElementById("deviceLocation").textContent = `${device.address}, ${device.city}`;

    // Update status
    updateDeviceStatus(device.device_state);

    // Update control buttons based on current state
    updateControlButtons(device.device_state);
  } catch (error) {
    console.error("Error loading device:", error);
    showMessage("Failed to load device information", "error");
  }
}

function updateDeviceStatus(stateId) {
  const statusEl = document.getElementById("deviceStatus");
  const states = {
    1: { name: "Active", class: "active" },
    2: { name: "Inactive", class: "inactive" },
    3: { name: "Maintenance", class: "maintenance" },
    4: { name: "Decommissioned", class: "inactive" }
  };

  const state = states[stateId] || { name: "Unknown", class: "inactive" };
  
  statusEl.className = `status-indicator ${state.class}`;
  statusEl.innerHTML = `
    <span class="status-dot"></span>
    <span>${state.name}</span>
  `;
}

function updateControlButtons(currentState) {
  const btnOn = document.getElementById("btnTurnOn");
  const btnOff = document.getElementById("btnTurnOff");
  const btnConfigure = document.getElementById("btnConfigure");

  // Enable/disable buttons based on current state
  if (currentState === 1) { // Active
    btnOn.disabled = true;
    btnOff.disabled = false;
    btnConfigure.disabled = false;
  } else if (currentState === 2) { // Inactive
    btnOn.disabled = false;
    btnOff.disabled = true;
    btnConfigure.disabled = false;
  } else {
    btnOn.disabled = true;
    btnOff.disabled = true;
    btnConfigure.disabled = true;
  }
}

async function controlDevice(action) {
  const token = localStorage.getItem("token");
  if (!token) {
    showMessage("Please login to control devices", "error");
    return;
  }

  try {
    const response = await fetch(`${API_BASE}/devices/${deviceId}/control`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ action }),
    });

    const data = await response.json();

    if (response.ok) {
      const actionNames = {
        4: "turned on",
        5: "turned off",
        6: "configured"
      };
      showMessage(`Device ${actionNames[action]} successfully. Current state: ${data.message.state}`, "success");
      
      // Reload device data and history
      setTimeout(() => {
        loadDeviceData();
        loadStateHistory();
      }, 500);
    } else {
      showMessage(data.error || "Control action failed", "error");
    }
  } catch (error) {
    console.error("Error controlling device:", error);
    showMessage("Failed to control device", "error");
  }
}

async function loadReadings() {
  try {
    const limit = document.getElementById("readingsLimit")?.value || 20;
    const response = await fetch(`${API_BASE}/devices/${deviceId}/readings?limit=${limit}`);
    
    if (!response.ok) throw new Error("Failed to load readings");

    const data = await response.json();
    const readings = data.readings || [];

    if (readings.length === 0) {
      document.getElementById("readingsTableBody").innerHTML = 
        '<tr><td colspan="4" style="text-align: center;">No readings available</td></tr>';
      return;
    }

    // Update live stats
    const latest = readings[0];
    if (latest) {
      document.getElementById("currentVoltage").textContent = `${latest.voltage.toFixed(2)} V`;
      document.getElementById("currentCurrent").textContent = `${latest.current.toFixed(2)} A`;
      document.getElementById("currentPower").textContent = `${(latest.voltage * latest.current).toFixed(2)} W`;
      document.getElementById("lastUpdate").textContent = new Date(latest.timestamp * 1000).toLocaleTimeString();
    }

    // Update table
    const tbody = document.getElementById("readingsTableBody");
    tbody.innerHTML = readings
      .map((r) => {
        const time = new Date(r.timestamp * 1000);
        const power = (r.voltage * r.current).toFixed(2);
        return `
          <tr>
            <td>${time.toLocaleString()}</td>
            <td>${r.voltage.toFixed(2)}</td>
            <td>${r.current.toFixed(2)}</td>
            <td>${power}</td>
          </tr>
        `;
      })
      .join("");

    // Update chart
    renderChart(readings);
  } catch (error) {
    console.error("Error loading readings:", error);
    document.getElementById("readingsTableBody").innerHTML = 
      '<tr><td colspan="4" style="text-align: center; color: red;">Failed to load readings</td></tr>';
  }
}

function renderChart(readings) {
  const voltageData = readings.reverse().map((r) => [r.timestamp * 1000, r.voltage]);
  const currentData = readings.map((r) => [r.timestamp * 1000, r.current]);

  if (voltageChart) {
    voltageChart.destroy();
  }

  voltageChart = Highcharts.chart("voltageChart", {
    chart: {
      type: "line",
      backgroundColor: "transparent",
    },
    title: {
      text: "Voltage & Current Over Time",
      style: { color: "#e5e7eb" },
    },
    xAxis: {
      type: "datetime",
      labels: { style: { color: "#9ca3af" } },
    },
    yAxis: [
      {
        title: { text: "Voltage (V)", style: { color: "#3b82f6" } },
        labels: { style: { color: "#9ca3af" } },
      },
      {
        title: { text: "Current (A)", style: { color: "#10b981" } },
        labels: { style: { color: "#9ca3af" } },
        opposite: true,
      },
    ],
    series: [
      {
        name: "Voltage",
        data: voltageData,
        color: "#3b82f6",
        yAxis: 0,
      },
      {
        name: "Current",
        data: currentData,
        color: "#10b981",
        yAxis: 1,
      },
    ],
    legend: {
      itemStyle: { color: "#e5e7eb" },
    },
    credits: { enabled: false },
  });
}

async function loadStateHistory() {
  try {
    // Note: This endpoint might need to be created in the backend
    // For now, we'll use a placeholder or adapt to existing endpoints
    const tbody = document.getElementById("historyTableBody");
    
    // Placeholder - you may need to create a proper endpoint
    tbody.innerHTML = `
      <tr>
        <td><span class="action-badge create">create</span></td>
        <td>Active</td>
        <td>${new Date(Date.now() - 86400000 * 3).toLocaleString()}</td>
        <td>admin</td>
      </tr>
      <tr>
        <td><span class="action-badge turn_off">turn_off</span></td>
        <td>Inactive</td>
        <td>${new Date(Date.now() - 86400000 * 2).toLocaleString()}</td>
        <td>admin</td>
      </tr>
      <tr>
        <td><span class="action-badge turn_on">turn_on</span></td>
        <td>Active</td>
        <td>${new Date(Date.now() - 86400000).toLocaleString()}</td>
        <td>admin</td>
      </tr>
    `;
  } catch (error) {
    console.error("Error loading state history:", error);
    document.getElementById("historyTableBody").innerHTML = 
      '<tr><td colspan="4" style="text-align: center; color: red;">Failed to load history</td></tr>';
  }
}

function showMessage(message, type = "info") {
  const msgEl = document.getElementById("controlMessage");
  msgEl.textContent = message;
  msgEl.style.display = "block";
  
  if (type === "success") {
    msgEl.style.background = "#d1fae5";
    msgEl.style.color = "#065f46";
    msgEl.style.border = "1px solid #10b981";
  } else if (type === "error") {
    msgEl.style.background = "#fee2e2";
    msgEl.style.color = "#991b1b";
    msgEl.style.border = "1px solid #ef4444";
  } else {
    msgEl.style.background = "#dbeafe";
    msgEl.style.color = "#1e40af";
    msgEl.style.border = "1px solid #3b82f6";
  }

  setTimeout(() => {
    msgEl.style.display = "none";
  }, 5000);
}

function exportToCSV() {
  const table = document.querySelector("#readingsTableBody");
  const rows = table.querySelectorAll("tr");
  
  let csv = "Timestamp,Voltage (V),Current (A),Power (W)\n";
  
  rows.forEach((row) => {
    const cols = row.querySelectorAll("td");
    if (cols.length === 4) {
      const rowData = Array.from(cols).map(col => col.textContent).join(",");
      csv += rowData + "\n";
    }
  });

  const blob = new Blob([csv], { type: "text/csv" });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `device-${deviceId}-readings-${Date.now()}.csv`;
  a.click();
  window.URL.revokeObjectURL(url);
}

// Cleanup on page unload
window.addEventListener("beforeunload", () => {
  if (autoRefreshInterval) {
    clearInterval(autoRefreshInterval);
  }
  if (voltageChart) {
    voltageChart.destroy();
  }
});
