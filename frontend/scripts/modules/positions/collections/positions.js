define([
  'backbone',
  '../models/position'
], function(Backbone, Position) {
  'use strict';

  return Backbone.Collection.extend({
    model: Position,
    comparator: 'invoiceDate'
  });
});
