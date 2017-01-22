'use strict';

var assert = require('assert'),
  test = require('selenium-webdriver/testing'),
  By = require('selenium-webdriver').By,
  until = require('selenium-webdriver').until,
  webdriver = require('selenium-webdriver');

test.describe('period management', function() {
  function buildDriver() {
    return new webdriver.Builder()
    .withCapabilities(webdriver.Capabilities.chrome())
    .build();
  }

  function buildFiscalPeriod(driver) {
    driver.get('http://umsatz.dev');
    driver.wait(until.elementLocated(By.css('.new-fiscalPeriod'))).then(function() {
      driver.findElement(By.css('.new-fiscalPeriod')).click();
      driver.findElement(By.css('#name')).sendKeys('2014 Q4');
      driver.findElement(By.css('#startsAt')).sendKeys('01/09/2014');
      driver.findElement(By.css('#endsAt')).sendKeys('01/12/2014');

      driver.findElement(By.css('.fiscal-period')).submit();
    });
  }

  function deleteFiscalPeriod(driver) {
    driver.get('http://umsatz.dev');
    driver.wait(until.elementLocated(By.css('table.fiscalPeriods'))).then(function() {
      driver.findElements(By.css('.delete-fiscalPeriod')).then(function(els) {
        els[els.length - 1].click();
      });
    });
  }

  test.it('create a new period', function() {
    var driver = buildDriver();

    buildFiscalPeriod(driver);

    driver.wait(until.elementLocated(By.css('.has-dropdown.fiscalPeriods'))).then(function() {
      driver.findElement(By.css('.has-dropdown.fiscalPeriods')).click();
      driver.findElements(By.css('.has-dropdown.fiscalPeriods .show-fiscalPeriod')).then(function(els) {
        els[els.length - 1].click();
      }).then(function() {
        driver.wait(until.elementLocated(By.css('.new-period'))).then(function() {
          driver.findElement(By.css('.new-period')).click().then(function() {
            driver.findElement(By.css('#accountCode1')).sendKeys('2000');
            driver.findElement(By.css('#accountLabel1')).sendKeys('Bank');
            driver.findElement(By.css('#accountCode2')).sendKeys('1000');
            driver.findElement(By.css('#accountLabel2')).sendKeys('Fremdleistungen');

            driver.findElement(By.css('#invoiceNumber')).sendKeys('20140101');
            driver.findElement(By.css('#totalAmount')).sendKeys('20.50');

            driver.findElement(By.css('.fiscal-item')).submit();
          });
        });
      });
      // driver.findElement(By.css('.new-fiscalPeriod')).click();
      // driver.findElement(By.css('#name')).sendKeys('2014 Q4');
      // driver.findElement(By.css('#startsAt')).sendKeys('01/09/2014');
      // driver.findElement(By.css('#endsAt')).sendKeys('01/12/2014');

      // driver.findElement(By.css('.fiscal-period')).submit()
      // .then(function() {
      //   driver.findElements(By.css('table.fiscalPeriods')).then(function(els) {
      //     els[0].getText().then(function(content) {
      //       assert.ok(content.match(/2014 Q4/));
      //     });
      //   });
      // });
    });

    deleteFiscalPeriod(driver);

    driver.quit();
  });

});