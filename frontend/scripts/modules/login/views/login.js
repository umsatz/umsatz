define([
  'application',
  'rv!../templates/login',
  'communicator'
], function(App, loginTemplate, comm) {
  'use strict';

  return function(registration) {
    var ractive = new App.Ractive({
      template: loginTemplate,
      adapt: ['Backbone'],

      el: '#content',
      data: {
        registration: registration,
        login: {}
      }
    });

    ractive.on('validate', function(context) {
      context.original.preventDefault();
      comm.reqres.request('session:validate', context.context.login.password, context.context.login.totp)
        .then(function() {
          App.Backbone.history.navigate('/', true);
        }).fail(function() {
          // TODO display errors
        });
    });

    return ractive;
  };
});
