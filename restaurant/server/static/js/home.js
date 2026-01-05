function toggleMenu() {
    const menu = document.getElementById('menu');
    const overlay = document.getElementById('overlay');
    const isMenuOpen = menu.classList.contains('menu-open');
    
    if (isMenuOpen) {
      menu.classList.remove('menu-open');
      overlay.classList.remove('overlay-active');
    } else {
      menu.classList.add('menu-open');
      overlay.classList.add('overlay-active');
    }
  }