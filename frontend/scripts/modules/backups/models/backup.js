define([
  'backbone'
], function(Backbone) {
  'use strict';

  return Backbone.Model.extend({
    url: function() {
      return '/api/backups' + (this.isNew() ? '' : '/' + this.get('id'));
    },

    link: function(rel) {
      var links = this.get('_links');
      for (var i = links.length - 1; i >= 0; i--) {
        var link = links[i];
        if (link.rel === rel) {
          return link;
        }
      }
      return null;
    },

    sync: function(method, model, options) {
      if (method.toLowerCase() === 'delete') {
        options = options || {};
        options.url = this.link('delete').href;
      }
      Backbone.sync(method, model, options);
    },

    restore: function() {
      var link = this.link('restore');
      return $.ajax({
        url: link.href,
        method: link.method
      });
    }
  });
});
