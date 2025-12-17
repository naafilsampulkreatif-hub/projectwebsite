import fs from 'fs/promises'
import path from 'path'
import pngToIco from 'png-to-ico'

async function run() {
  const publicDir = path.resolve(new URL(import.meta.url).pathname, '..', '..', 'public')
  // On Windows the pathname may start with a slash; normalize
  const normalizedPublic = publicDir.replace(/^\/+/, '')
  const input = path.join(normalizedPublic, 'NS-logo.png')
  const output = path.join(normalizedPublic, 'favicon.ico')

  try {
    const buf = await fs.readFile(input)
    const ico = await pngToIco(buf)
    await fs.writeFile(output, ico)
    console.log('favicon.ico written to', output)
  } catch (err) {
    console.error('Failed to generate favicon.ico:', err.message || err)
    process.exit(1)
  }
}

run()
