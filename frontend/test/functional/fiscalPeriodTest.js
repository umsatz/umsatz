'use strict';

var assert = require('assert'),
  test = require('selenium-webdriver/testing'),
  By = require('selenium-webdriver').By,
  until = require('selenium-webdriver').until,
  webdriver = require('selenium-webdriver');

test.describe('fiscalPeriod management', function() {
  function buildDriver() {
    return new webdriver.Builder()
    .withCapabilities(webdriver.Capabilities.chrome())
    .build();
  }

  test.it('create a new fiscalPeriod', function() {
    var driver = buildDriver();

    driver.get('http://umsatz.dev');
    driver.wait(until.elementLocated(By.css('.new-fiscalPeriod'))).then(function() {
      driver.findElement(By.css('.new-fiscalPeriod')).click();
      driver.findElement(By.css('#name')).sendKeys('2014 Q4');
      driver.findElement(By.css('#startsAt')).sendKeys('01/09/2014');
      driver.findElement(By.css('#endsAt')).sendKeys('01/12/2014');

      driver.findElement(By.css('.fiscal-period')).submit()
      .then(function() {
        driver.findElements(By.css('table.fiscalPeriods')).then(function(els) {
          els[0].getText().then(function(content) {
            assert.ok(content.match(/2014 Q4/));
          });
        });
      });
    });

    driver.quit();
  });

  test.it('update existing fiscalPeriod', function() {
    var driver = buildDriver();

    driver.get('http://umsatz.dev');
    driver.wait(until.elementLocated(By.css('table.fiscalPeriods'))).then(function() {
      driver.findElements(By.css('table.fiscalPeriods .edit-fiscalPeriod')).then(function(els) {
        els[els.length - 1].click();
      }).then(function() {
        driver.findElement(By.css('#name')).sendKeys('2015 Q1');
        driver.findElement(By.css('#startsAt')).sendKeys('01/01/2015');
        driver.findElement(By.css('#endsAt')).sendKeys('01/04/2015');

        driver.findElement(By.css('.fiscal-period')).submit();
      }).then(function() {
        driver.findElements(By.css('table.fiscalPeriods')).then(function(els) {
          els[0].getText().then(function(content) {
            assert.ok(content.match(/2015 Q1/));
          });
        });
      });
    });

    driver.quit();
  });

  test.it('delete fiscalPeriod', function() {
    var driver = buildDriver();

    driver.get('http://umsatz.dev');
    driver.wait(until.elementLocated(By.css('table.fiscalPeriods'))).then(function() {
      driver.findElements(By.css('.delete-fiscalPeriod')).then(function(els) {
        els[els.length - 1].click();
      }).then(function() {
        driver.findElements(By.css('table.fiscalPeriods')).then(function(els) {
          els[0].getText().then(function(content) {
            assert.ok(content.match(/2015 Q1/) === null);
          });
        });
      });
    });

    driver.quit();
  });
});