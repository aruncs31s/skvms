const tokenKey = "skvms_token";
const userKey = "skvms_user";

const getToken = () => localStorage.getItem(tokenKey);
const setToken = (token) => localStorage.setItem(tokenKey, token);
const clearToken = () => localStorage.removeItem(tokenKey);

const getUser = () => {
  const user = localStorage.getItem(userKey);
  return user ? JSON.parse(user) : null;
};
const setUser = (user) => localStorage.setItem(userKey, JSON.stringify(user));
const clearUser = () => localStorage.removeItem(userKey);

// Load navbar
(async () => {
  const container = document.getElementById("navbar-container");
  if (container) {
    const res = await fetch("/static/navbar.html");
    const html = await res.text();
    container.innerHTML = html;
    initNavbar();
  }
})();

function initNavbar() {
  const navAuthBadge = document.getElementById("navAuthBadge");
  const navLogoutBtn = document.getElementById("navLogoutBtn");
  const navLoginLink = document.getElementById("navLoginLink");
  const navManageDevices = document.getElementById("navManageDevices");
  const navManageUsers = document.getElementById("navManageUsers");
  const navAudit = document.getElementById("navAudit");

  if (navLogoutBtn) {
    navLogoutBtn.addEventListener("click", () => {
      clearToken();
      clearUser();
      window.location.href = "/";
    });
  }

  updateAuthUI();

  function updateAuthUI() {
    const user = getUser();
    const displayText = user ? `Welcome, ${user.name}` : "Guest";
    
    if (navAuthBadge) navAuthBadge.textContent = displayText;
    if (navLoginLink) navLoginLink.style.display = user ? "none" : "inline-block";
    if (navLogoutBtn) navLogoutBtn.style.display = user ? "inline-block" : "none";
    
    // Show admin links only when logged in
    if (navManageDevices) navManageDevices.style.display = user ? "inline-block" : "none";
    if (navManageUsers) navManageUsers.style.display = user ? "inline-block" : "none";
    if (navAudit) navAudit.style.display = user ? "inline-block" : "none";

    // Mark active page
    const path = window.location.pathname;
    document.querySelectorAll(".navbar-item").forEach((item) => {
      const href = item.getAttribute("href");
      if (href && (href === path || (href !== "/" && path.startsWith(href)))) {
        item.classList.add("active");
      }
    });
  }
}



const authBadge = document.getElementById("authBadge");
const logoutBtn = document.getElementById("logoutBtn");
const loginLink = document.getElementById("loginLink");

if (logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    clearToken();
    clearUser();
    window.location.reload();
  });
}

const updateAuthUI = () => {
  const user = getUser();
  const displayText = user ? `Welcome, ${user.name}` : "Guest";
  
  if (authBadge) {
    authBadge.textContent = displayText;
  }
  if (loginLink) {
    loginLink.style.display = user ? "none" : "inline-block";
  }
  if (logoutBtn) {
    logoutBtn.style.display = user ? "inline-block" : "none";
  }
};

updateAuthUI();

const getDeviceIdFromPath = () => {
  const match = window.location.pathname.match(/^\/devices\/(\d+)$/);
  if (!match) return null;
  return match[1];
};

const loginForm = document.getElementById("loginForm");
if (loginForm) {
  loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const formData = new FormData(loginForm);
    const payload = {
      username: formData.get("username"),
      password: formData.get("password"),
    };

    const messageEl = document.getElementById("loginMessage");
    messageEl.textContent = "";
    try {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      if (!res.ok) {
        messageEl.textContent = data.error || "Login failed";
        return;
      }
      setToken(data.token);
      setUser(data.user);
      updateAuthUI();
      messageEl.textContent = "Login successful. Redirecting...";
      setTimeout(() => {
        window.location.href = "/";
      }, 700);
    } catch (error) {
      messageEl.textContent = "Login failed. Try again.";
    }
  });
}

const deviceList = document.getElementById("deviceList");
if (deviceList) {
  loadDevices();
}

const deviceId = getDeviceIdFromPath();
if (deviceId) {
  loadDevicePage(deviceId);
}

const allReadingsPage = window.location.pathname === "/all-readings";
if (allReadingsPage) {
  loadAllReadingsPage();
}

const manageDevicesPage = window.location.pathname === "/manage-devices";
if (manageDevicesPage) {
  loadManageDevicesPage();
}

const manageUsersPage = window.location.pathname === "/manage-users";
if (manageUsersPage) {
  loadManageUsersPage();
}

const auditPage = window.location.pathname === "/audit";
if (auditPage) {
  loadAuditPage();
}

async function loadDevices() {
  deviceList.innerHTML = "<p class=\"muted\">Loading devices...</p>";
  try {
    const res = await fetch("/api/devices");
    const data = await res.json();
    if (!res.ok) {
      deviceList.innerHTML = "<p class=\"muted\">Failed to load devices.</p>";
      return;
    }
    renderDevices(data.devices || []);
  } catch (error) {
    deviceList.innerHTML = "<p class=\"muted\">Failed to load devices.</p>";
  }
}

function renderDevices(devices) {
  if (!devices.length) {
    deviceList.innerHTML = "<p class=\"muted\">No devices found.</p>";
    return;
  }

  const token = getToken();
  deviceList.innerHTML = devices
    .map(
      (device) => `
      <div class="device-card clickable" data-device-id="${device.id}">
        <div class="device-header">
          <h3>${device.name}</h3>
          <span class="chip">${device.type}</span>
        </div>
        <div class="device-body">
          <p><strong>IP:</strong> ${device.ip_address || "-"}</p>
          <p><strong>MAC:</strong> ${device.mac_address || "-"}</p>
          <p><strong>Firmware:</strong> ${device.firmware_version || "-"}</p>
          <p><strong>Location:</strong> ${device.address || "-"}, ${device.city || "-"}</p>
        </div>
        <div class="device-actions">
          <button class="button" ${token ? "" : "disabled"} data-id="${device.id}">Send Control</button>
          <a class="button ghost" href="/devices/${device.id}">View Readings</a>
          <span class="hint">${token ? "" : "Login required to control"}</span>
        </div>
      </div>
    `
    )
    .join("");

  deviceList.querySelectorAll(".device-card.clickable").forEach((card) => {
    card.addEventListener("click", (e) => {
      const target = e.target;
      if (target && (target.tagName === "BUTTON" || target.tagName === "A")) {
        return;
      }
      const id = card.getAttribute("data-device-id");
      if (id) window.location.href = `/devices/${id}`;
    });
  });

  deviceList.querySelectorAll("button[data-id]").forEach((button) => {
    button.addEventListener("click", async () => {
      const id = button.getAttribute("data-id");
      await sendControlCommand(id);
    });
  });
}

async function loadDevicePage(deviceId) {
  const titleEl = document.getElementById("deviceTitle");
  const metaEl = document.getElementById("deviceMeta");
  const readingsBody = document.getElementById("readingsBody");
  const emptyEl = document.getElementById("readingsEmpty");
  const limitEl = document.getElementById("historyLimit");
  const refreshBtn = document.getElementById("refreshBtn");
  const datePickerEl = document.getElementById("datePicker");
  const showYesterdayEl = document.getElementById("showYesterday");

  const voltageEl = document.getElementById("latestVoltage");
  const currentEl = document.getElementById("latestCurrent");
  const timeEl = document.getElementById("latestTime");
  const maxVoltageEl = document.getElementById("maxVoltage");
  const minVoltageEl = document.getElementById("minVoltage");
  const maxVoltageTimeEl = document.getElementById("maxVoltageTime");
  const minVoltageTimeEl = document.getElementById("minVoltageTime");

  let chart = null;

  const fmtTime = (unixSeconds) => {
    if (!unixSeconds) return "--";
    const d = new Date(unixSeconds * 1000);
    return d.toLocaleString();
  };

  const fmtNum = (value, digits) => {
    if (value === null || value === undefined) return "--";
    if (Number.isNaN(Number(value))) return "--";
    return Number(value).toFixed(digits);
  };

  const fetchDevice = async () => {
    const res = await fetch(`/api/devices/${deviceId}`);
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "failed to load device");
    return data.device;
  };

  const fetchReadings = async (limit) => {
    const res = await fetch(`/api/devices/${deviceId}/readings?limit=${limit}`);
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "failed to load readings");
    return data;
  };

  const fetchReadingsByDate = async (startTime, endTime) => {
    const res = await fetch(`/api/devices/${deviceId}/readings/range?start=${startTime}&end=${endTime}`);
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "failed to load readings");
    return data;
  };

  const renderChart = (readings, yesterdayReadings = []) => {
    const voltageData = readings.map((r) => [r.timestamp * 1000, r.voltage]);
    const currentData = readings.map((r) => [r.timestamp * 1000, r.current]);
    
    const series = [
      {
        name: "Voltage (V)",
        data: voltageData,
        color: "#2563eb",
        yAxis: 0,
      },
      {
        name: "Current (A)",
        data: currentData,
        color: "#10b981",
        yAxis: 1,
      },
    ];

    if (yesterdayReadings.length > 0) {
      const yVoltageData = yesterdayReadings.map((r) => [r.timestamp * 1000, r.voltage]);
      const yCurrentData = yesterdayReadings.map((r) => [r.timestamp * 1000, r.current]);
      series.push({
        name: "Yesterday Voltage (V)",
        data: yVoltageData,
        color: "#93c5fd",
        dashStyle: "ShortDash",
        yAxis: 0,
      });
      series.push({
        name: "Yesterday Current (A)",
        data: yCurrentData,
        color: "#6ee7b7",
        dashStyle: "ShortDash",
        yAxis: 1,
      });
    }

    if (chart) {
      chart.destroy();
    }

    chart = Highcharts.chart("chartContainer", {
      chart: {
        type: "spline",
        backgroundColor: "#ffffff",
      },
      title: {
        text: "Device Readings Over Time",
        style: { fontWeight: "700", fontSize: "18px" },
      },
      xAxis: {
        type: "datetime",
        title: { text: "Time" },
      },
      yAxis: [
        {
          title: { text: "Voltage (V)" },
          labels: { format: "{value} V" },
        },
        {
          title: { text: "Current (A)" },
          labels: { format: "{value} A" },
          opposite: true,
        },
      ],
      tooltip: {
        shared: true,
        crosshairs: true,
      },
      legend: {
        enabled: true,
      },
      series: series,
    });
  };

  const renderReadings = (payload, statsData = null) => {
    const latest = payload.latest;
    const readings = payload.readings || [];

    if (latest) {
      voltageEl.textContent = `${fmtNum(latest.voltage, 1)} V`;
      currentEl.textContent = `${fmtNum(latest.current, 2)} A`;
      timeEl.textContent = fmtTime(latest.timestamp);
    } else {
      voltageEl.textContent = "--";
      currentEl.textContent = "--";
      timeEl.textContent = "--";
    }

    if (statsData) {
      maxVoltageEl.textContent = `${fmtNum(statsData.max_voltage, 1)} V`;
      minVoltageEl.textContent = `${fmtNum(statsData.min_voltage, 1)} V`;
      maxVoltageTimeEl.textContent = fmtTime(statsData.max_voltage_time);
      minVoltageTimeEl.textContent = fmtTime(statsData.min_voltage_time);
    }

    readingsBody.innerHTML = "";
    if (!readings.length) {
      emptyEl.style.display = "block";
      return;
    }
    emptyEl.style.display = "none";
    readingsBody.innerHTML = readings
      .map(
        (r) => `
        <tr>
          <td>${fmtTime(r.timestamp)}</td>
          <td>${fmtNum(r.voltage, 1)}</td>
          <td>${fmtNum(r.current, 2)}</td>
        </tr>
      `
      )
      .join("");
  };

  const refresh = async () => {
    try {
      if (metaEl) metaEl.textContent = "Loading readings...";

      const device = await fetchDevice();
      if (titleEl) titleEl.textContent = device?.name || `Device ${deviceId}`;
      if (metaEl)
        metaEl.textContent = `${device?.type || ""} • IP ${device?.ip_address || "-"} • MAC ${device?.mac_address || "-"}`;

      const selectedDate = datePickerEl ? datePickerEl.value : null;
      const showYesterday = showYesterdayEl ? showYesterdayEl.checked : false;

      let readingsPayload;
      let yesterdayPayload = null;

      if (selectedDate) {
        const start = new Date(selectedDate);
        start.setHours(0, 0, 0, 0);
        const end = new Date(selectedDate);
        end.setHours(23, 59, 59, 999);
        
        readingsPayload = await fetchReadingsByDate(Math.floor(start.getTime() / 1000), Math.floor(end.getTime() / 1000));
        
        if (showYesterday) {
          const yStart = new Date(start);
          yStart.setDate(yStart.getDate() - 1);
          const yEnd = new Date(end);
          yEnd.setDate(yEnd.getDate() - 1);
          yesterdayPayload = await fetchReadingsByDate(Math.floor(yStart.getTime() / 1000), Math.floor(yEnd.getTime() / 1000));
        }
      } else {
        const limit = limitEl ? Number(limitEl.value) : 50;
        readingsPayload = await fetchReadings(limit);
        
        if (showYesterday) {
          const now = new Date();
          const yStart = new Date(now);
          yStart.setDate(yStart.getDate() - 1);
          yStart.setHours(0, 0, 0, 0);
          const yEnd = new Date(now);
          yEnd.setDate(yEnd.getDate() - 1);
          yEnd.setHours(23, 59, 59, 999);
          yesterdayPayload = await fetchReadingsByDate(Math.floor(yStart.getTime() / 1000), Math.floor(yEnd.getTime() / 1000));
        }
      }

      renderReadings(readingsPayload, readingsPayload.stats);
      renderChart(
        readingsPayload.readings || [],
        yesterdayPayload ? yesterdayPayload.readings || [] : []
      );
    } catch (error) {
      if (metaEl) metaEl.textContent = "Failed to load device/readings.";
      if (emptyEl) emptyEl.style.display = "block";
    }
  };

  if (limitEl) limitEl.addEventListener("change", refresh);
  if (refreshBtn) refreshBtn.addEventListener("click", refresh);
  if (datePickerEl) datePickerEl.addEventListener("change", refresh);
  if (showYesterdayEl) showYesterdayEl.addEventListener("change", refresh);

  const exportCsvBtn = document.getElementById("exportCsvBtn");
  const exportJsonBtn = document.getElementById("exportJsonBtn");

  if (exportCsvBtn) {
    exportCsvBtn.addEventListener("click", () => {
      const selectedDate = datePickerEl ? datePickerEl.value : null;
      if (selectedDate) {
        const start = new Date(selectedDate);
        start.setHours(0, 0, 0, 0);
        const end = new Date(selectedDate);
        end.setHours(23, 59, 59, 999);
        exportReadingsCSV(deviceId, Math.floor(start.getTime() / 1000), Math.floor(end.getTime() / 1000));
      } else {
        const limit = limitEl ? Number(limitEl.value) : 50;
        exportReadingsCSV(deviceId, null, null, limit);
      }
    });
  }

  if (exportJsonBtn) {
    exportJsonBtn.addEventListener("click", () => {
      const selectedDate = datePickerEl ? datePickerEl.value : null;
      if (selectedDate) {
        const start = new Date(selectedDate);
        start.setHours(0, 0, 0, 0);
        const end = new Date(selectedDate);
        end.setHours(23, 59, 59, 999);
        exportReadingsJSON(deviceId, Math.floor(start.getTime() / 1000), Math.floor(end.getTime() / 1000));
      } else {
        const limit = limitEl ? Number(limitEl.value) : 50;
        exportReadingsJSON(deviceId, null, null, limit);
      }
    });
  }

  await refresh();
}

async function exportReadingsCSV(deviceId, startTime, endTime, limit) {
  try {
    let url = `/api/devices/${deviceId}/readings`;
    if (startTime && endTime) {
      url = `/api/devices/${deviceId}/readings/range?start=${startTime}&end=${endTime}`;
    } else if (limit) {
      url += `?limit=${limit}`;
    }

    const res = await fetch(url);
    const data = await res.json();
    const readings = data.readings || [];

    if (!readings.length) {
      alert("No readings to export");
      return;
    }

    let csv = "Time,Voltage (V),Current (A)\n";
    readings.forEach((r) => {
      const time = new Date(r.timestamp * 1000).toLocaleString();
      csv += `"${time}",${r.voltage},${r.current}\n`;
    });

    const blob = new Blob([csv], { type: "text/csv" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `device-${deviceId}-readings.csv`;
    link.click();
  } catch (error) {
    alert("Export failed");
  }
}

async function exportReadingsJSON(deviceId, startTime, endTime, limit) {
  try {
    let url = `/api/devices/${deviceId}/readings`;
    if (startTime && endTime) {
      url = `/api/devices/${deviceId}/readings/range?start=${startTime}&end=${endTime}`;
    } else if (limit) {
      url += `?limit=${limit}`;
    }

    const res = await fetch(url);
    const data = await res.json();

    const blob = new Blob([JSON.stringify(data, null, 2)], {
      type: "application/json",
    });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `device-${deviceId}-readings.json`;
    link.click();
  } catch (error) {
    alert("Export failed");
  }
}

async function loadAllReadingsPage() {
  const allDevicesGrid = document.getElementById("allDevicesGrid");
  const refreshAllBtn = document.getElementById("refreshAllBtn");
  const exportAllCsvBtn = document.getElementById("exportAllCsvBtn");
  let combinedChart = null;

  const fetchAllDevices = async () => {
    const res = await fetch("/api/devices");
    const data = await res.json();
    if (!res.ok) throw new Error("failed to load devices");
    return data.devices || [];
  };

  const fetchDeviceReadings = async (deviceId, limit = 20) => {
    const res = await fetch(`/api/devices/${deviceId}/readings?limit=${limit}`);
    const data = await res.json();
    if (!res.ok) return { latest: null, readings: [] };
    return data;
  };

  const renderAllDevices = async () => {
    try {
      allDevicesGrid.innerHTML = "<p class='muted'>Loading...</p>";
      const devices = await fetchAllDevices();

      const devicesData = await Promise.all(
        devices.map(async (device) => {
          const readingsData = await fetchDeviceReadings(device.id);
          return { device, ...readingsData };
        })
      );

      allDevicesGrid.innerHTML = devicesData
        .map(
          ({ device, latest }) => `
        <div class="device-reading-card">
          <div class="device-reading-header">
            <h3>${device.name}</h3>
            <span class="chip">${device.type}</span>
          </div>
          <div class="device-reading-stats">
            <div class="mini-stat">
              <div class="mini-stat-label">Voltage</div>
              <div class="mini-stat-value">${latest ? latest.voltage.toFixed(1) : "--"} V</div>
            </div>
            <div class="mini-stat">
              <div class="mini-stat-label">Current</div>
              <div class="mini-stat-value">${latest ? latest.current.toFixed(2) : "--"} A</div>
            </div>
          </div>
          <div class="device-reading-footer">
            <a href="/devices/${device.id}" class="button ghost">View Details</a>
          </div>
        </div>
      `
        )
        .join("");

      renderCombinedChart(devicesData);
    } catch (error) {
      allDevicesGrid.innerHTML = "<p class='muted'>Failed to load devices.</p>";
    }
  };

  const renderCombinedChart = (devicesData) => {
    const series = [];
    const colors = ["#2563eb", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"];

    devicesData.forEach((item, idx) => {
      const readings = item.readings || [];
      if (readings.length > 0) {
        const voltageData = readings.map((r) => [r.timestamp * 1000, r.voltage]);
        series.push({
          name: `${item.device.name} - Voltage`,
          data: voltageData,
          color: colors[idx % colors.length],
          yAxis: 0,
        });
      }
    });

    if (combinedChart) {
      combinedChart.destroy();
    }

    combinedChart = Highcharts.chart("combinedChart", {
      chart: { type: "spline", backgroundColor: "#ffffff" },
      title: { text: "All Devices - Voltage Comparison", style: { fontWeight: "700", fontSize: "18px" } },
      xAxis: { type: "datetime", title: { text: "Time" } },
      yAxis: { title: { text: "Voltage (V)" }, labels: { format: "{value} V" } },
      tooltip: { shared: true, crosshairs: true },
      legend: { enabled: true },
      series: series,
    });
  };

  const exportAllCSV = async () => {
    try {
      const devices = await fetchAllDevices();
      let csv = "Device,Time,Voltage (V),Current (A)\n";

      for (const device of devices) {
        const readingsData = await fetchDeviceReadings(device.id, 50);
        const readings = readingsData.readings || [];
        readings.forEach((r) => {
          const time = new Date(r.timestamp * 1000).toLocaleString();
          csv += `"${device.name}","${time}",${r.voltage},${r.current}\n`;
        });
      }

      const blob = new Blob([csv], { type: "text/csv" });
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "all-devices-readings.csv";
      link.click();
    } catch (error) {
      alert("Export failed");
    }
  };

  if (refreshAllBtn) refreshAllBtn.addEventListener("click", renderAllDevices);
  if (exportAllCsvBtn) exportAllCsvBtn.addEventListener("click", exportAllCSV);

  await renderAllDevices();
}

async function loadManageDevicesPage() {
  const devicesTableBody = document.getElementById("devicesTableBody");
  const addDeviceBtn = document.getElementById("addDeviceBtn");
  const deviceModal = document.getElementById("deviceModal");
  const closeModal = document.getElementById("closeModal");
  const deviceForm = document.getElementById("deviceForm");

  let currentDeviceId = null;
  let deviceTypes = [];

  const fetchDeviceTypes = async () => {
    const res = await fetch("/api/device-types");
    const data = await res.json();
    return data || [];
  };

  const populateDeviceTypeDropdown = () => {
    const select = document.getElementById("deviceType");
    select.innerHTML = '<option value="">Select device type</option>';
    deviceTypes.forEach(type => {
      const option = document.createElement("option");
      option.value = type.id;
      option.textContent = type.name;
      select.appendChild(option);
    });
  };

  const fetchDevices = async () => {
    const res = await fetch("/api/devices");
    const data = await res.json();
    return data.devices || [];
  };

  const renderDevices = async () => {
    const devices = await fetchDevices();
    devicesTableBody.innerHTML = devices
      .map(
        (d) => `
      <tr>
        <td>${d.id}</td>
        <td>${d.name}</td>
        <td>${d.type}</td>
        <td>${d.ip_address || "-"}</td>
        <td>${d.mac_address || "-"}</td>
        <td>
          <button class="button ghost" data-edit="${d.id}" data-device='${JSON.stringify(d)}'>Edit</button>
          <button class="button ghost" data-delete="${d.id}">Delete</button>
        </td>
      </tr>
    `
      )
      .join("");

    devicesTableBody.querySelectorAll("[data-edit]").forEach((btn) => {
      btn.addEventListener("click", () => {
        const device = JSON.parse(btn.getAttribute("data-device"));
        openEditModal(device);
      });
    });
    devicesTableBody.querySelectorAll("[data-delete]").forEach((btn) => {
      btn.addEventListener("click", () => deleteDevice(btn.getAttribute("data-delete")));
    });
  };

  const openEditModal = (device) => {
    if (device) {
      document.getElementById("modalTitle").textContent = "Edit Device";
      currentDeviceId = device.id;
      document.getElementById("deviceName").value = device.name;
      // Find the device type ID by name
      const typeOption = deviceTypes.find(t => t.name === device.type);
      document.getElementById("deviceType").value = typeOption ? typeOption.id : "";
      document.getElementById("deviceIP").value = device.ip_address || "";
      document.getElementById("deviceMAC").value = device.mac_address || "";
      document.getElementById("deviceFirmware").value = device.firmware_version || "";
      document.getElementById("deviceAddress").value = device.address || "";
      document.getElementById("deviceCity").value = device.city || "";
    } else {
      document.getElementById("modalTitle").textContent = "Add Device";
      currentDeviceId = null;
      deviceForm.reset();
    }
    deviceModal.style.display = "flex";
  };

  addDeviceBtn.addEventListener("click", () => openEditModal(null));
  closeModal.addEventListener("click", () => (deviceModal.style.display = "none"));

  deviceForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const token = getToken();
    if (!token) {
      alert("Please login first");
      return;
    }

    const payload = {
      name: document.getElementById("deviceName").value,
      type: parseInt(document.getElementById("deviceType").value),
      ip_address: document.getElementById("deviceIP").value,
      mac_address: document.getElementById("deviceMAC").value,
      firmware_version: document.getElementById("deviceFirmware").value,
      address: document.getElementById("deviceAddress").value,
      city: document.getElementById("deviceCity").value,
    };

    try {
      const url = currentDeviceId ? `/api/devices/${currentDeviceId}` : "/api/devices";
      const method = currentDeviceId ? "PUT" : "POST";
      const res = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      if (!res.ok) {
        alert(data.error || "Failed to save device");
        return;
      }

      alert(data.message);
      deviceModal.style.display = "none";
      await renderDevices();
    } catch (error) {
      alert("Failed to save device");
    }
  });

  const deleteDevice = async (id) => {
    if (!confirm("Delete this device?")) return;
    const token = getToken();
    if (!token) {
      alert("Please login first");
      return;
    }

    try {
      const res = await fetch(`/api/devices/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await res.json();
      if (!res.ok) {
        alert(data.error || "Failed to delete device");
        return;
      }

      alert(data.message);
      await renderDevices();
    } catch (error) {
      alert("Failed to delete device");
    }
  };

  // Load device types and then render
  try {
    deviceTypes = await fetchDeviceTypes();
    populateDeviceTypeDropdown();
    await renderDevices();
  } catch (error) {
    console.error("Failed to load device types:", error);
    await renderDevices();
  }
}

async function loadManageUsersPage() {
  const usersTableBody = document.getElementById("usersTableBody");
  const addUserBtn = document.getElementById("addUserBtn");
  const userModal = document.getElementById("userModal");
  const closeModal = document.getElementById("closeModal");
  const userForm = document.getElementById("userForm");

  let currentUserId = null;

  const fetchUsers = async () => {
    const token = getToken();
    const res = await fetch("/api/users", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    const data = await res.json();
    return data.users || [];
  };

  const renderUsers = async () => {
    try {
      const users = await fetchUsers();
      usersTableBody.innerHTML = users
        .map(
          (u) => `
        <tr>
          <td>${u.id}</td>
          <td>${u.username}</td>
          <td>${u.email}</td>
          <td>${u.role}</td>
          <td>
            <button class="button ghost" data-edit="${u.id}" data-user='${JSON.stringify(u)}'>Edit</button>
            <button class="button ghost" data-delete="${u.id}">Delete</button>
          </td>
        </tr>
      `
        )
        .join("");

      usersTableBody.querySelectorAll("[data-edit]").forEach((btn) => {
        btn.addEventListener("click", () => {
          const user = JSON.parse(btn.getAttribute("data-user"));
          openEditModal(user);
        });
      });
      usersTableBody.querySelectorAll("[data-delete]").forEach((btn) => {
        btn.addEventListener("click", () => deleteUser(btn.getAttribute("data-delete")));
      });
    } catch (error) {
      usersTableBody.innerHTML = `<tr><td colspan="5" class="muted">Failed to load users. Please login.</td></tr>`;
    }
  };

  const openEditModal = (user) => {
    if (user) {
      document.getElementById("modalTitle").textContent = "Edit User";
      currentUserId = user.id;
      document.getElementById("username").value = user.username;
      document.getElementById("email").value = user.email;
      document.getElementById("password").value = "";
      document.getElementById("password").placeholder = "Leave blank to keep current password";
      document.getElementById("role").value = user.role;
    } else {
      document.getElementById("modalTitle").textContent = "Add User";
      currentUserId = null;
      userForm.reset();
      document.getElementById("password").placeholder = "Password";
    }
    userModal.style.display = "flex";
  };

  addUserBtn.addEventListener("click", () => openEditModal(null));
  closeModal.addEventListener("click", () => (userModal.style.display = "none"));

  userForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const token = getToken();
    if (!token) {
      alert("Please login first");
      return;
    }

    const payload = {
      username: document.getElementById("username").value,
      email: document.getElementById("email").value,
      password: document.getElementById("password").value,
      role: document.getElementById("role").value,
    };

    // For update, password is optional
    if (currentUserId && !payload.password) {
      delete payload.password;
    }

    try {
      const url = currentUserId ? `/api/users/${currentUserId}` : "/api/users";
      const method = currentUserId ? "PUT" : "POST";
      const res = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      if (!res.ok) {
        alert(data.error || "Failed to save user");
        return;
      }

      alert(data.message);
      userModal.style.display = "none";
      await renderUsers();
    } catch (error) {
      alert("Failed to save user");
    }
  });

  const deleteUser = async (id) => {
    if (!confirm("Delete this user?")) return;
    const token = getToken();
    if (!token) {
      alert("Please login first");
      return;
    }

    try {
      const res = await fetch(`/api/users/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await res.json();
      if (!res.ok) {
        alert(data.error || "Failed to delete user");
        return;
      }

      alert(data.message);
      await renderUsers();
    } catch (error) {
      alert("Failed to delete user");
    }
  };

  await renderUsers();
}


async function loadAuditPage() {
  const auditTableBody = document.getElementById("auditTableBody");
  const auditEmpty = document.getElementById("auditEmpty");
  const actionFilter = document.getElementById("actionFilter");
  const refreshAuditBtn = document.getElementById("refreshAuditBtn");
  const exportAuditBtn = document.getElementById("exportAuditBtn");

  const fetchAuditLogs = async (action = "") => {
    const url = action ? `/api/audit?action=${action}` : "/api/audit";
    const res = await fetch(url);
    const data = await res.json();
    return data.logs || [];
  };

  const renderAuditLogs = async () => {
    try {
      const action = actionFilter ? actionFilter.value : "";
      const logs = await fetchAuditLogs(action);

      if (!logs.length) {
        auditTableBody.innerHTML = "";
        auditEmpty.style.display = "block";
        return;
      }

      auditEmpty.style.display = "none";
      auditTableBody.innerHTML = logs
        .map(
          (log) => `
        <tr>
          <td>${new Date(log.created_at).toLocaleString()}</td>
          <td>${log.username}</td>
          <td><span class="chip">${log.action}</span></td>
          <td>${log.details}</td>
          <td>${log.ip_address}</td>
        </tr>
      `
        )
        .join("");
    } catch (error) {
      auditTableBody.innerHTML = `<tr><td colspan="5" class="muted">Failed to load audit logs</td></tr>`;
    }
  };

  const exportAuditCSV = async () => {
    try {
      const logs = await fetchAuditLogs(actionFilter ? actionFilter.value : "");
      let csv = "Time,User,Action,Details,IP Address\n";
      logs.forEach((log) => {
        const time = new Date(log.created_at).toLocaleString();
        csv += `"${time}","${log.username}","${log.action}","${log.details}","${log.ip_address}"\n`;
      });

      const blob = new Blob([csv], { type: "text/csv" });
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "audit-log.csv";
      link.click();
    } catch (error) {
      alert("Export failed");
    }
  };

  if (actionFilter) actionFilter.addEventListener("change", renderAuditLogs);
  if (refreshAuditBtn) refreshAuditBtn.addEventListener("click", renderAuditLogs);
  if (exportAuditBtn) exportAuditBtn.addEventListener("click", exportAuditCSV);

  await renderAuditLogs();
}

async function sendControlCommand(id) {
  const token = getToken();
  if (!token) {
    alert("Login required to control devices.");
    return;
  }

  try {
    const res = await fetch(`/api/devices/${id}/control`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ command: "connect" }),
    });
    const data = await res.json();
    if (!res.ok) {
      alert(data.error || "Command failed");
      return;
    }
    alert(data.message || "Command sent");
  } catch (error) {
    alert("Command failed");
  }
}