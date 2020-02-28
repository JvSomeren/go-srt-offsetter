export const debounce = (fn, delay) => {
  let timeout;

  return function() {
    clearTimeout(timeout);

    timeout = setTimeout(() => {
      fn.apply(this, arguments);
    }, delay);
  }
};

export const languages = {
  'en': 'English',
  'nl': 'Dutch',
};
