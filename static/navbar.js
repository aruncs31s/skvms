// navbar.js - Navbar functionality

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