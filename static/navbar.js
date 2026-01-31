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
  const navLogoutBtn = document.getElementById("navLogoutBtn");

  if (navLogoutBtn) {
    navLogoutBtn.addEventListener("click", () => {
      clearToken();
      clearUser();
      window.location.href = "/";
    });
  }

  updateAuthUI();
}