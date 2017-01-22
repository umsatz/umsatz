define([
  'backbone'
], function(Backbone) {
  'use strict';

  var Position = Backbone.Model.extend({
    defaults: function() {
      var now = new Date();
      return {
        type: 'expense',
        invoiceDate: now.getFullYear() + '-' + ((now.getMonth() + 1) < 10 ? '0' : '') + (now.getMonth() + 1) + '-' + (now.getDate() < 10 ? '0' : '') + now.getDate(),
        invoiceNumber: '',
        totalAmountCents: 0,
        currency: 'EUR',
        tax: 19,
        description: '',
        attachment: null,
        accountCodeFrom: '',
        accountCodeTo: '',
        errors: {}
      };
    },

    url: function() {
      return '/api/positions/' + (this.isNew() ? '' : '' + this.get('id'));
    },

    hasErrorOn: function(attr) {
      return attr in this.get('errors');
    },

    invoiceMonth: function() {
      return this.get('invoiceDate').split('-')[1];
    },

    isIncome: function() {
      return this.get('type') === 'income';
    },

    totalVatAmountCents: function() {
      return this.get('totalAmountCents') - (this.get('totalAmountCents') / (this.get('tax') / 100 + 1));
    },

    signedTotalAmountCents: function() {
      if (this.get('type') === 'expense') {
        return this.get('totalAmountCents') * -1;
      }
      return this.get('totalAmountCents');
    },

    toJSON: function() {
      var data = Backbone.Model.prototype.toJSON.apply(this);
      if (data !== null) {
        data.accountCodeFrom = data.accountCodeFrom.toString();
        data.accountCodeTo = data.accountCodeTo.toString();
        data.description = data.description.toString();
        data.invoiceNumber = data.invoiceNumber.toString();
        data.tax = parseInt(data.tax, 10);
      }
      return data;
    }
  });
  return Position;
});
