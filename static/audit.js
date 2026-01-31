// audit.js - Audit page functionality

const auditPage = window.location.pathname === "/audit";
if (auditPage) {
  loadAuditPage();
}

async function loadAuditPage() {
  const auditTableBody = document.getElementById("auditTableBody");
  if (!auditTableBody) return;

  auditTableBody.innerHTML = "<tr><td colspan=\"5\">Loading audit logs...</td></tr>";

  try {
    const res = await fetch("/api/audit", {
      headers: {
        "Authorization": `Bearer ${getToken()}`
      }
    });
    const data = await res.json();

    if (!res.ok) {
      auditTableBody.innerHTML = "<tr><td colspan=\"5\">Failed to load audit logs</td></tr>";
      return;
    }

    renderAuditTable(data || []);
  } catch (error) {
    auditTableBody.innerHTML = "<tr><td colspan=\"5\">Failed to load audit logs</td></tr>";
  }

  // Export CSV button
  const exportCsvBtn = document.getElementById("exportAuditCsvBtn");
  if (exportCsvBtn) {
    exportCsvBtn.addEventListener("click", exportAuditCSV);
  }
}

function renderAuditTable(auditLogs) {
  const auditTableBody = document.getElementById("auditTableBody");

  if (!auditLogs.length) {
    auditTableBody.innerHTML = "<tr><td colspan=\"5\">No audit logs found</td></tr>";
    return;
  }

  auditTableBody.innerHTML = auditLogs.map(log => `
    <tr>
      <td>${new Date(log.created_at).toLocaleString()}</td>
      <td>${log.username}</td>
      <td><span class="chip">${log.action}</span></td>
      <td>${log.details}</td>
      <td>${log.ip_address}</td>
    </tr>
  `).join('');
}

function exportAuditCSV() {
  const auditTableBody = document.getElementById("auditTableBody");
  if (!auditTableBody) return;

  const rows = auditTableBody.querySelectorAll("tr");
  if (!rows.length) return;

  let csv = "Time,Username,Action,Details,IP Address\n";

  rows.forEach(row => {
    const cells = row.querySelectorAll("td");
    if (cells.length === 5) {
      const time = new Date(cells[0].textContent).toLocaleString();
      const username = cells[1].textContent;
      const action = cells[2].textContent.replace('chip', '').trim();
      const details = cells[3].textContent;
      const ip = cells[4].textContent;

      csv += `"${time}","${username}","${action}","${details}","${ip}"\n`;
    }
  });

  const blob = new Blob([csv], { type: 'text/csv' });
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `audit-log-${new Date().toISOString().split('T')[0]}.csv`;
  link.click();
  window.URL.revokeObjectURL(url);
}