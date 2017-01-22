define([
  'application',
  'rv!templates/footer',
  'i18n',
  'jquery',
  'foundation.dropdown',
], function(App, FooterTemplate, I18n, jQuery) {
  'use strict';

  var footer = new App.Ractive({
    template: FooterTemplate,

    el: '[role=footer]',

    complete: function() {
      jQuery(document).foundation();
    },

    selectLanguage: function(language) {
      I18n.setLng(language, function() {
        window.location.reload();
      });
    }
  });
  return footer;
});
