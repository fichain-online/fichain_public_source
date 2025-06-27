'use client';

import { useEffect, useMemo, useState, useCallback, type ReactNode } from 'react';
import Particles, { initParticlesEngine } from '@tsparticles/react';
import { loadSlim } from '@tsparticles/slim';

export function ParticlesWrapper({ children }: { children: ReactNode }) {
  const [init, setInit] = useState(false);

  // This should be run only once per application lifetime
  useEffect(() => {
    initParticlesEngine(async (engine) => {
      // You can initiate the tsParticles instance (engine) here, adding custom shapes or presets
      // This loads the slim version of tsParticles, which is enough for most cases
      await loadSlim(engine);
    }).then(() => {
      setInit(true);
    });
  }, []);

  const particlesLoaded = useCallback(async (container: any) => {
    // You can perform actions once the particles are loaded
    // console.log('Particles container loaded', container);
  }, []);

  const options = useMemo(() => ({
    fpsLimit: 60,
    interactivity: {
      events: {
        onClick: {
          enable: true,
          mode: "push",
        },
        onHover: {
          enable: true,
          mode: "repulse",
        },
      },
      modes: {
        push: {
          quantity: 2,
        },
        repulse: {
          distance: 100,
          duration: 0.4,
        },
      },
    },
    particles: {
      color: {
        value: "#a0aec0", // A light gray (slate-400)
      },
      links: {
        color: "#cbd5e0", // A slightly darker gray for links (slate-300)
        distance: 150,
        enable: true,
        opacity: 0.2,
        width: 1,
      },
      collisions: {
        enable: true,
      },
      move: {
        direction: "none",
        enable: true,
        outModes: {
          default: "bounce",
        },
        random: false,
        speed: 0.5,
        straight: false,
      },
      number: {
        density: {
          enable: true,
          area: 1000,
        },
        value: 40,
      },
      opacity: {
        value: 0.3,
      },
      shape: {
        type: "circle",
      },
      size: {
        value: { min: 1, max: 3 },
      },
    },
    detectRetina: true,
    background: {
      color: {
        value: "transparent", // This is crucial for your page background to show through
      },
    },
  }), []);

  return (
    <>
      {init && (
        <Particles
          id="tsparticles"
          particlesLoaded={particlesLoaded}
          options={options}
          className="fixed top-0 left-0 w-full h-full z-[-1]"
        />
      )}
      {children}
    </>
  );
}
