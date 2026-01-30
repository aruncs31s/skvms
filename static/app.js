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
  if (authBadge) {
    authBadge.textContent = user ? `Welcome, ${user.name}` : "Guest";
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

  await refresh();
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