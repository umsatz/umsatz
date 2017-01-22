define([
  'hbs/handlebars',
  'i18n'
], function(Handlebars, I18n) {
  'use strict';

  function t(key, options) {
    if (options === null) {
      options = {};
    } else {
      options = options.hash;
    }

    return I18n.t(key, options);
  }

  Handlebars.registerHelper('t', t);
  return t;
});
