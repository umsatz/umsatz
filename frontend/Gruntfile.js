'use strict';
var LIVERELOAD_PORT = 35729;
var SERVER_PORT = 9000;
var lrSnippet = require('connect-livereload')({
  port: LIVERELOAD_PORT
});
var mountFolder = function(connect, dir) {
  return connect.static(require('path').resolve(dir));
};
module.exports = function(grunt) {
  // load all grunt tasks
  require('load-grunt-tasks')(grunt);
  // show elapsed time at the end
  require('time-grunt')(grunt);
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-connect');
  grunt.loadNpmTasks('grunt-contrib-cssmin');
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-requirejs');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-file-blocks');
  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-webdriver');

  // configurable paths
  var yeomanConfig = {
    app: '',
    dist: 'dist'
  };

  grunt.initConfig({
    yeoman: yeomanConfig,

    // watch list
    watch: {
      options: {
        nospawn: true,
        livereload: true
      },
      livereload: {
        files: [
          'scripts/**/*.js',
          'templates/{,**/}*.hbs',
          'test/spec/**/*.js',
          'index.html'
        ]
      },
      sass: {
        files: 'styles/**/*.scss',
        tasks: ['sass']
      },
      // compass: {
      //     files: ['styles/**/*.scss'],
      //     tasks: ['compass']
      // },
    },

    webdriver: {
      options: {
        desiredCapabilities: {
          browserName: 'chrome'
        }
      },
      sometests: {
        tests: ['test/functional/*.js']
      }
    },

    // testing server
    connect: {
      options: {
        port: SERVER_PORT,
        hostname: 'localhost'
      },
      livereload: {
        options: {
          middleware: function(connect) {
            return [
              lrSnippet,
              mountFolder(connect, '.tmp'),
              mountFolder(connect, yeomanConfig.app),
              mountFolder(connect, 'bower_components/foundation-icon-fonts')
            ];
          }
        }
      },
      dist: {
        options: {
          middleware: function(connect) {
            return [
              mountFolder(connect, yeomanConfig.dist)
            ];
          }
        }
      }
    },

    sass: {
      options: {
        includePaths: [
          'bower_components/foundation/scss',
          'bower_components/foundation-icon-fonts',
          'bower_components'
        ]
      },
      dist: {
        options: {
          outputStyle: 'compressed'
        },
        files: {
          '.tmp/styles/main.css': 'styles/main.scss'
        }
      }
    },

    clean: {
      dist: ['.tmp', '<%= yeoman.dist %>/*'],
      server: '.tmp'
    },

    // linting
    jshint: {
      options: {
        jshintrc: '.jshintrc',
        reporter: require('jshint-stylish')
      },
      all: [
        'Gruntfile.js',
        'scripts/**/*.js',
        '!scripts/vendor/*',
        'test/spec/{,*/}*.js'
      ]
    },

    requirejs: {
      dist: {
        options: {
          baseUrl: 'scripts',
          mainConfigFile: 'scripts/init.js',
          include: 'main',
          name: '../bower_components/almond/almond',
          out: 'dist/scripts/main.js',
          wrap: true,
          findNestedDependencies: true
        }
      }
    },

    copy: {
      dist: {
        files: [{
            expand: true,
            dot: true,
            cwd: '',
            dest: '<%= yeoman.dist %>',
            src: [
              '*.{ico,txt}',
              '.htaccess',
              'index.html',
              'images/{,*/}*.{webp,gif}',
              'styles/fonts/{,*/}*.*'
            ]
          }, {
            expand: true,
            dot: true,
            cwd: '',
            dest: '<%= yeoman.dist %>',
            src: ['fixtures/**']
          }, {
            expand: true,
            dot: true,
            cwd: '',
            dest: '<%= yeoman.dist %>',
            src: ['locales/*.json']
          }, {
            expand: true,
            dot: true,
            cwd: 'bower_components/foundation-icon-fonts/',
            dest: '<%= yeoman.dist %>',
            src: [
              '*.{ttf,eot,woff,svg}'
            ]
          }, {
            expand: true,
            dot: true,
            cwd: 'bower_components/pdfjs-dist/build/',
            dest: '<%= yeoman.dist %>/bower_components/pdfjs-dist/build/',
            src: ['*.js']
          }]
      }
    },

    cssmin: {
      dist: {
        files: {
          '<%= yeoman.dist %>/styles/main.css': [
            '.tmp/styles/**/*.css',
            'styles/**/*.css'
          ]
        }
      }
    },

    fileblocks: {
      options: {
        templates: {
          'js': '<script data-main="scripts/init" src="${file}"></script>',
        },
        removeFiles: true
      },
      prod: {
        src: 'dist/index.html',
        blocks: {
          'app': {
            src: 'scripts/main.js'
          }
        }
      }
    },
  });

  grunt.registerTask('createDefaultTemplate', function() {
    grunt.file.write('.tmp/scripts/templates.js', 'this.JST = this.JST || {};');
  });

  // starts express server with live testing via testserver
  grunt.registerTask('default', function() {
    grunt.task.run([
      'clean:server',
      'sass',
      'connect:livereload',
      'watch'
    ]);
  });

  grunt.registerTask('test', function() {
    grunt.task.run([
      'webdriver'
    ]);
  });

  grunt.registerTask('build', [
    'clean',
    'createDefaultTemplate',
    'sass',
    'requirejs',
    'cssmin',
    'copy',
    'fileblocks:prod'
  ]);

};
