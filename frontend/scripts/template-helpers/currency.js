define(['accounting'], function(accounting) {
  'use strict';

  var symbols = {
    'EUR': '€',
    'USD': '$'
  };

  function currency(amount, options) {
    if (options === null || options === undefined) {
      options = {};
    }
    options.currency = options.currency || 'EUR';
    options.cents = options.cents || 100;

    amount = parseFloat(amount) || 0;
    amount /= options.cents;

    return accounting.formatMoney(amount, {
      symbol: symbols[options.currency],
      format: '%v %s',
      thousand: '.',
      decimal: ','
    });
  }

  return currency;
});
