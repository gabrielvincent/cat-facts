window.loadLazyComponent = function loadLazyComponent(lazyId, elementId) {
  const suspenseEl = document.querySelector(`div[data-lazy-id="${lazyId}"]`);
  const lazyEl = document.querySelector(`#${elementId}`);

  console.log(`--- will load:`, lazyId, elementId);
  if (suspenseEl == null || lazyEl == null) {
    console.log(`--- null:`, suspenseEl, lazyEl);
    return;
  }

  suspenseEl.innerHTML = lazyEl.innerHTML;
  lazyEl.remove();
};
