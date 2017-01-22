/* global PDFJS */
'use strict';
define([
  'pdfviewer'
], function() {
  PDFJS.workerSrc = '../bower_components/pdfjs-dist/build/pdf.worker.js';

  var previewPDF = function($el, path) {
    if (!(PDFJS.PDFViewer && PDFJS.getDocument)) {
      return;
    }
    var container = $el[0];
    var pdfViewer = new PDFJS.PDFViewer({
      container: container
    });

    container.addEventListener('pagesinit', function () {
      // we can use pdfViewer now, e.g. let's change default scale.
      pdfViewer.currentScaleValue = 'page-width';
    });

    PDFJS.getDocument(path).then(function (pdfDocument) {
      pdfViewer.setDocument(pdfDocument);
    });
  };

  var previewJPG = function($el, path) {
    var image = $('<img src="' + path + '"/>');
    $el.append(image);
  };

  var previewMapping = {
    'pdf': previewPDF,
    'jpg': previewJPG
  };

  return {
    canPreview: function(path) {
      return path !== null && path !== undefined && previewMapping[path.toLowerCase().split('.').pop()] != null;
    },
    display: function($el, path) {
      var fileType = path.toLowerCase().split('.').pop();
      previewMapping[fileType].apply(this, [$el, path]);
    }
  };
});
