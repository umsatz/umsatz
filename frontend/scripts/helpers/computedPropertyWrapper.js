define([], function() {
  'use strict';
  return function(attr, obj) {
    return {
      get: function() {
        return obj[attr];
      },
      set: function(value) {
        obj[attr] = value;
      }
    };
  };
});
