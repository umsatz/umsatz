define([
  'i18next'
], function(I18n) {
  'use strict';

  I18n.init({
    interpolationPrefix: '%{',
    interpolationSuffix: '}',
    resGetPath: '../locales/%{lng}.json',
    preload: ['de', 'en'],
    fallbackLng: 'de',
    useCookie: true,
  });

  return I18n;
});
