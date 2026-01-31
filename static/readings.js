// readings.js - All readings page functionality

const allReadingsPage = window.location.pathname === "/all-readings";
if (allReadingsPage) {
  loadAllReadingsPage();
  const refreshAllBtn = document.getElementById("refreshAllBtn");
  if (refreshAllBtn) {
    refreshAllBtn.addEventListener("click", () => {
      loadAllReadingsPage();
    });
  }
}

async function loadAllReadingsPage() {
  console.log("Loading all readings page...");
  const readingsContainer = document.getElementById("allReadingsContainer");
  if (!readingsContainer) {
    console.log("allReadingsContainer not found!");
    return;
  }

  readingsContainer.innerHTML = "<p class=\"muted\">Loading all readings...</p>";

  try {
    const res = await fetch("/api/devices");
    const data = await res.json();
    console.log("Devices API response:", data);

    if (!res.ok) {
      readingsContainer.innerHTML = "<p class=\"muted\">Failed to load devices.</p>";
      return;
    }

    const devices = data.devices || [];
    console.log("Devices found:", devices);
    if (!devices.length) {
      readingsContainer.innerHTML = "<p class=\"muted\">No devices found.</p>";
      return;
    }

    // Load readings for all devices
    const allReadings = [];
    console.log("Loading readings for", devices.length, "devices");
    for (const device of devices) {
      try {
        console.log("Loading readings for device:", device.id, device.name);
        const readingsRes = await fetch(`/api/devices/${device.id}/readings?limit=10`);
        const readingsData = await readingsRes.json();
        console.log("Readings for device", device.id, ":", readingsData);
        if (readingsRes.ok && readingsData.readings) {
          readingsData.readings.forEach(reading => {
            allReadings.push({
              ...reading,
              deviceName: device.name,
              deviceId: device.id
            });
          });
        }
      } catch (error) {
        console.error(`Failed to load readings for device ${device.id}:`, error);
      }
    }

    console.log("Total readings collected:", allReadings.length);

    // Sort by timestamp descending
    allReadings.sort((a, b) => b.timestamp - a.timestamp);

    renderAllReadings(allReadings);
  } catch (error) {
    readingsContainer.innerHTML = "<p class=\"muted\">Failed to load readings.</p>";
  }
}

function renderAllReadings(readings) {
  console.log("Rendering readings:", readings);
  const readingsContainer = document.getElementById("allReadingsContainer");

  if (!readings.length) {
    readingsContainer.innerHTML = "<p class=\"muted\">No readings found.</p>";
    return;
  }

  readingsContainer.innerHTML = `
    <div class="table-wrap">
      <table class="table">
        <thead>
          <tr>
            <th>Device</th>
            <th>Time</th>
            <th>Voltage (V)</th>
            <th>Current (A)</th>
            <th>Power (W)</th>
          </tr>
        </thead>
        <tbody>
          ${readings.map(reading => `
            <tr>
              <td><a href="/devices/${reading.deviceId}" class="table-link">${reading.deviceName}</a></td>
              <td>${new Date(reading.timestamp * 1000).toLocaleString()}</td>
              <td>${reading.voltage.toFixed(2)}</td>
              <td>${reading.current.toFixed(2)}</td>
              <td>${(reading.voltage * reading.current).toFixed(2)}</td>
            </tr>
          `).join('')}
        </tbody>
      </table>
    </div>
  `;

  // Create combined chart
  const deviceGroups = {};
  readings.forEach(reading => {
    if (!deviceGroups[reading.deviceId]) {
      deviceGroups[reading.deviceId] = {
        name: reading.deviceName,
        data: []
      };
    }
    deviceGroups[reading.deviceId].data.push([
      reading.timestamp * 1000,
      reading.voltage
    ]);
  });

  const series = Object.values(deviceGroups).map(group => ({
    name: group.name,
    data: group.data.sort((a, b) => a[0] - b[0])
  }));

  console.log("Chart series:", series);

  if (typeof Highcharts !== 'undefined') {
    Highcharts.chart('combinedChart', {
      title: { text: 'Combined Device Voltages' },
      xAxis: { type: 'datetime' },
      yAxis: { title: { text: 'Voltage (V)' } },
      series: series
    });
  } else {
    console.error("Highcharts not loaded");
  }
}