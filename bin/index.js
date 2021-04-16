#!/usr/bin/env node
const { Command } = require('commander');
const chalk = require('chalk');
const clipboardy = require('clipboardy');

// next packages:
require('@jswork/next');
require('@jswork/next-absolute-package');

const { version } = nx.absolutePackage();
const program = new Command();

program.version(version);

program
  .option('-n, --npm-install', 'Get npm install script.')
  .option('-s, --shortname', 'Get short name.')
  .parse(process.argv);

nx.declare({
  statics: {
    init() {
      const app = new this();
      app.start();
    }
  },
  methods: {
    init() {},
    start() {
      if (program.shortname) {
        const [_, shortname] = nx.absolutePackage('name').split('/');
        clipboardy.writeSync(shortname);
      }

      if (program.npmInstall) {
        const name = nx.absolutePackage('name');
        clipboardy.writeSync(`npm i ${name}`);
      }
    }
  }
});
