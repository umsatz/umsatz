define([
  'application',
  'communicator',
  'rv!../templates/preview'
], function(App, comm, PreviewTemplate) {
  'use strict';

  function loadPreview(id) {
    return $.get(window.location.protocol + '//' + window.location.hostname + '/api/reporting/preview/' + id);
  }

  return function(fiscalPeriod, template) {
    var ractive = new App.Ractive({
      template: PreviewTemplate,
      adapt: ['Backbone'],

      el: '#content',

      data: {
        fiscalYear: fiscalPeriod,
        template: template.get('data'),
        preview: '',
      }
    });

    ractive.observe('template', function(currentValue) {
      template.set('data', currentValue);
    });

    ractive.on('download', function(ctx) {
      ctx.original.preventDefault();

      var xhr = new XMLHttpRequest();
      xhr.open('GET', '/api/reporting/generate/' + fiscalPeriod.get('id'));
      // xhr.responseType = 'arraybuffer';
      xhr.responseType = 'blob';
      xhr.setRequestHeader('X-UMSATZ-SESSION', comm.reqres.request('session:current').key);

      xhr.onload = function () {
        if (this.status === 200) {
          var blob = xhr.response;
          var objectUrl = URL.createObjectURL(blob);
          window.open(objectUrl);
        }
      };
      xhr.send();
    });

    var iframeRefresh = function() {
      template.save().then(function() {
        var iframe = $('iframe', ractive.el);

        loadPreview(fiscalPeriod.get('id')).then(function(data) {
          iframe.attr('srcdoc', data);
        });
      });
    };

    var debouncedIframeRefresh = _.debounce(iframeRefresh.bind(this), 250);
    template.on('change', debouncedIframeRefresh);

    loadPreview(fiscalPeriod.get('id')).then(function(data) {
      ractive.set('preview', data);
    });

    return ractive;
  };
});
