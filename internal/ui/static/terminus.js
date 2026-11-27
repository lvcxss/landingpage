/**
 * Terminus CSS Helper JS
 * Provides accessible event delegation for interactive tabbed subtimelines
 * and terminal components without requiring external dependencies.
 */
(function () {
  'use strict';

  function switchSubtimelineTab(button) {
    const container = button.closest('.terminal-subtimeline');
    if (!container) return;

    const targetSelector = button.getAttribute('data-target');
    if (!targetSelector) return;

    // Set active tab button
    const buttons = container.querySelectorAll('.subtimeline-tab-btn');
    buttons.forEach((btn) => btn.classList.remove('active'));
    button.classList.add('active');

    // Switch panels
    const panels = container.querySelectorAll('.subtimeline-panel');
    panels.forEach((p) => p.classList.remove('active'));

    const targetPanel = container.querySelector(targetSelector);
    if (targetPanel) {
      targetPanel.classList.add('active');
    }
  }

  // Global click delegation for any subtimeline tabs
  document.addEventListener('click', function (event) {
    const tabBtn = event.target.closest('.subtimeline-tab-btn');
    if (tabBtn) {
      event.preventDefault();
      switchSubtimelineTab(tabBtn);
    }
  });

  // Expose function globally if needed
  window.switchSubtimelineTab = switchSubtimelineTab;
})();
