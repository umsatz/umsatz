'use strict';
define([
  'application',
  'communicator',
  './views/login',
  './views/signup'
], function(App, Communicator, loginView, signupView) {
  var LoginModule = App.module('LoginModule');

  var SessionService = (function() {
    var activeRegistration = null;
    var activeSession = null;
    if (localStorage.umsatzSession !== undefined) {
      activeSession = JSON.parse(localStorage.umsatzSession);
    }

    return {
      registration: function() {
        if (activeRegistration === null) {
          activeRegistration = new jQuery.Deferred();
          jQuery.get('/api/auth/registration/').then(function(registration) {
            activeRegistration.resolve(JSON.parse(registration));
          }).fail(function(error) {
            activeRegistration.reject(error);
          });
        }
        return activeRegistration;
      },
      session: function() {
        return activeSession;
      },
      validate: function(password, totp) {
        var p = new jQuery.Deferred();
        jQuery.post('/api/auth/signin/', JSON.stringify({
          password: password,
          otp: totp
        })).then(function(session) {
          activeSession = JSON.parse(session);
          localStorage.umsatzSession = session;
          Communicator.trigger('user:signin');
          p.resolve(activeSession);
        }).fail(function(error) {
          p.reject(error);
        });
        return p;
      }
    };
  })();

  Communicator.reqres.reply('session:validate', SessionService.validate.bind(SessionService));
  Communicator.reqres.reply('session:current', SessionService.session.bind(SessionService));
  Communicator.reqres.reply('registration', SessionService.registration.bind(SessionService));

  var Router = App.Backbone.Marionette.AppRouter.extend({
    appRoutes: {
      'login': 'loginOrSignup',
    },
    controller: {
      loginOrSignup: function() {
        SessionService.registration().then(function(reg) {
          loginView(reg);
        }).fail(function() {
          signupView();
        });
      }
    }
  });

  App.addInitializer(function() {
    new Router();
  });

  return LoginModule;
});
