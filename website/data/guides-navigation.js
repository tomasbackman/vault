export default [
  {
    category: 'getting-started'
  },
  {
    category: 'operations',
    content: [
      'reference-architecture',
      'vault-ha-consul',
      'production',
      'generate-root',
      'rekeying-and-rotating',
      'plugin-backends',
      '--------------',
      'replication',
      'disaster-recovery',
      'mount-filter',
      'multi-tenant',
      'autounseal-aws-kms',
      'seal-wrap',
      'monitoring'
    ]
  },
  {
    category: 'identity',
    content: [
      'secure-intro',
      'policies',
      'authentication',
      'approle-trusted-entities',
      'lease',
      'identity',
      '--------------',
      'sentinel',
      'control-groups'
    ]
  },
  {
    category: 'secret-mgmt',
    content: [
      'static-secrets',
      'versioned-kv',
      'dynamic-secrets',
      'db-root-rotation',
      'cubbyhole',
      'ssh-otp',
      'pki-engine',
      'app-integration'
    ]
  },
  {
    category: 'encryption',
    content: ['transit', 'spring-demo', 'transit-rewrap']
  }
]
