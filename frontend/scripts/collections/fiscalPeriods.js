define([
  'backbone',
  'models/fiscalPeriod'
], function(Backbone, FiscalPeriod) {
  'use strict';

  return Backbone.Collection.extend({
    model: FiscalPeriod,
    url: '/api/fiscalPeriods/',
    comparator: function(a, b) {
      return a.get('startsAt') < b.get('startsAt');
    }
  });
});
