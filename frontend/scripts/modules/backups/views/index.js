/* global $, alert */
define([
  'application',
  'rv!../templates/index'
], function(App, BackupsTemplate) {
  'use strict';

  return function(backups) {
    var ractive = new App.Ractive({
      template: BackupsTemplate,
      adapt: ['Backbone'],

      el: 'content',

      data: {
        backups: backups
      }
    });

    ractive.on('create', function(event) {
      $.ajax({
        url: '/api/backups',
        method: 'POST'
      }).then(function() {
        alert('success');
      }, function() {
        alert('error');
      });

      event.original.preventDefault();
      return false;
    });

    ractive.on('delete', function(event) {
      event.context.destroy();
    });

    ractive.on('restore', function(event) {
      event.context.restore();
    });

    return ractive;
  };
});
