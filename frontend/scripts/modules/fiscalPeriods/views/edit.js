define([
  'application',
  'rv!../templates/edit',
], function(App, PeriodTemplate) {
  'use strict';

  return function(period) {
    var originalAttributes = _.clone(period.attributes);

    var ractive = new App.Ractive({
      template: PeriodTemplate,
      adapt: ['Backbone'],

      el: 'content',

      data: {
        period: period,
      },
    });

    ractive.on({
      save: function(evt) {
        evt.original.preventDefault();

        period.save()
        .then(function(data) {
          period.set(data);

          this.fire('fiscalPeriod:put', period);
        }.bind(this));
      },
      cancel: function(evt) {
        evt.original.preventDefault();
        period.set(originalAttributes, {
          silence: true
        });
        this.fire('fiscalPeriod:cancel');
      }
    });

    return ractive;
  };
});
