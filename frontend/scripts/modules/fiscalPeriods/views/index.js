define([
  'application',
  'rv!../templates/index',
], function(App, Template) {
  'use strict';

  return function(fiscalPeriods) {
    var ractive = new App.Ractive({
      template: Template,
      adapt: ['Backbone'],

      el: 'content',

      data: {
        fiscalPeriods: fiscalPeriods
      },
    });

    ractive.on('delete', function(event) {
      event.context.destroy();
    });

    return ractive;
  };
});
