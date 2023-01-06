
/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
}

const signals = [
  "SIGHUP",
  "SIGINT",
  "SIGQUIT",
  "SIGILL",
  "SIGTRAP",
  "SIGABRT",
  "SIGBUS",
  "SIGFPE",
  "SIGUSR1",
  "SIGSEGV",
  "SIGUSR2",
  "SIGTERM",
];

const terminator = (signal) => {
  console.log(`Received ${signal}: ", "cleaning up`);
  setTimeout(
    (...args) => {
      console.log(...args);
      process.exit(0);
    },
    400,
    "Message"
  );
  // un-commenting this line will cause everything to just stop and kill the pending timeout
  // process.exit(0);
};

signals.map((signal) => {
  process.on(signal, () => terminator(signal));
});

module.exports = nextConfig
