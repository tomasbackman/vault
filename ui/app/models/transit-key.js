import { alias } from '@ember/object/computed';
import { set, get, computed } from '@ember/object';
import DS from 'ember-data';
import clamp from 'vault/utils/clamp';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const { attr } = DS;

const ACTION_VALUES = {
  encrypt: {
    isSupported: 'supportsEncryption',
    description: 'Looks up wrapping properties for the given token',
    glyph: 'lock-closed',
  },
  decrypt: {
    isSupported: 'supportsDecryption',
    description: 'Decrypts the provided ciphertext using this key',
    glyph: 'envelope-unsealed--outline',
  },
  datakey: {
    isSupported: 'supportsEncryption',
    description: 'Generates a new key and value encrypted with this key',
    glyph: 'key',
  },
  rewrap: {
    isSupported: 'supportsEncryption',
    description: 'Rewraps the ciphertext using the latest version of the named key',
    glyph: 'refresh-default',
  },
  sign: {
    isSupported: 'supportsSigning',
    description: 'Get the cryptographic signature of the given data',
    glyph: 'edit',
  },
  hmac: { isSupported: true, description: 'Generate a data digest using a hash algorithm', glyph: 'remix' },
  verify: {
    isSupported: true,
    description: 'Validate the provided signature for the given data',
    glyph: 'check-circle-outline',
  },
  export: { isSupported: 'exportable', description: 'Get the named key', glyph: 'exit' },
};

export default DS.Model.extend({
  type: attr('string', {
    defaultValue: 'aes256-gcm96',
  }),
  name: attr('string'),
  deletionAllowed: attr('boolean'),
  derived: attr('boolean'),
  exportable: attr('boolean'),
  minDecryptionVersion: attr('number', {
    defaultValue: 1,
  }),
  minEncryptionVersion: attr('number', {
    defaultValue: 0,
  }),
  latestVersion: attr('number'),
  keys: attr('object'),
  convergentEncryption: attr('boolean'),
  convergentEncryptionVersion: attr('number'),

  supportsSigning: attr('boolean'),
  supportsEncryption: attr('boolean'),
  supportsDecryption: attr('boolean'),
  supportsDerivation: attr('boolean'),

  setConvergentEncryption(val) {
    if (val === true) {
      set(this, 'derived', val);
    }
    set(this, 'convergentEncryption', val);
  },

  setDerived(val) {
    if (val === false) {
      set(this, 'convergentEncryption', val);
    }
    set(this, 'derived', val);
  },

  supportedActions: computed('type', function() {
    return Object.keys(ACTION_VALUES)
      .filter(name => {
        const { isSupported } = ACTION_VALUES[name];
        return typeof isSupported === 'boolean' || get(this, isSupported);
      })
      .map(name => {
        const { description, glyph } = ACTION_VALUES[name];
        return { name, description, glyph };
      });
  }),

  canDelete: computed('deletionAllowed', 'lastLoadTS', function() {
    const deleteAttrChanged = Boolean(this.changedAttributes().deletionAllowed);
    return get(this, 'deletionAllowed') && deleteAttrChanged === false;
  }),

  keyVersions: computed('validKeyVersions', function() {
    let maxVersion = Math.max(...get(this, 'validKeyVersions'));
    let versions = [];
    while (maxVersion > 0) {
      versions.unshift(maxVersion);
      maxVersion--;
    }
    return versions;
  }),

  encryptionKeyVersions: computed('keyVerisons', 'minDecryptionVersion', 'latestVersion', function() {
    const { keyVersions, minDecryptionVersion } = this.getProperties('keyVersions', 'minDecryptionVersion');

    return keyVersions
      .filter(version => {
        return version >= minDecryptionVersion;
      })
      .reverse();
  }),

  keysForEncryption: computed('minEncryptionVersion', 'latestVersion', function() {
    let { minEncryptionVersion, latestVersion } = this.getProperties('minEncryptionVersion', 'latestVersion');
    let minVersion = clamp(minEncryptionVersion - 1, 0, latestVersion);
    let versions = [];
    while (latestVersion > minVersion) {
      versions.push(latestVersion);
      latestVersion--;
    }
    return versions;
  }),

  validKeyVersions: computed('keys', function() {
    return Object.keys(get(this, 'keys'));
  }),

  exportKeyTypes: computed('exportable', 'type', function() {
    let types = ['hmac'];
    if (this.get('supportsSigning')) {
      types.unshift('signing');
    }
    if (this.get('supportsEncryption')) {
      types.unshift('encryption');
    }
    return types;
  }),

  backend: attr('string'),

  rotatePath: lazyCapabilities(apiPath`${'backend'}/keys/${'id'}/rotate`, 'backend', 'id'),
  canRotate: alias('rotatePath.canUpdate'),
  secretPath: lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id'),
  canRead: alias('secretPath.canUpdate'),
  canEdit: alias('secretPath.canUpdate'),
});
