define([
  'backbone'
], function(Backbone) {
  'use strict';

  return Backbone.Model.extend({
    url: function() {
      return '/api/fiscalPeriods/' + (this.isNew() ? '' : '' + this.get('id'));
    },
  });
});
