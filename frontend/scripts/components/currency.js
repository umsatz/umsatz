define(
['ractive', 'template-helpers/currency'], function(Ractive, currencyHelper) {
  'use strict';

  var CurrencyWidget = Ractive.extend({
    template: '<span class="{{ type }}">{{ amount }}</span>',
    init: function() {
      this.setAmount();
    },
    setAmount: function() {
      if (this.data.amountCents !== null) {
        this.observe('amountCents', function(newValue) {
          this.set('amount', currencyHelper(newValue, {
            cents: 100,
            currency: this.data.currency
          }));
        });
        this.set('amount', currencyHelper(this.data.amountCents, {
          cents: 100,
          currency: this.data.currency
        }));
      } else {
        this.set('amount', currencyHelper(this.data.amount, {
          cents: 1,
          currency: this.data.currency
        }));
      }
    },
    data: {
      amount: 'missing amount'
    }
  });
  Ractive.components.currency = CurrencyWidget;
});
