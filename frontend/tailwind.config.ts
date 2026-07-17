import type { Config } from 'tailwindcss';
const { heroui } = require('@heroui/react');

const withVar = (name: string) => `rgb(var(${name}) / <alpha-value>)`;

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
    './node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        paper: withVar('--paper'),
        'paper-dim': withVar('--paper-dim'),
        ink: withVar('--ink'),
        'ink-soft': withVar('--ink-soft'),
        'ink-faint': withVar('--ink-faint'),
        line: withVar('--line'),
        card: withVar('--card'),
        sidebar: withVar('--sidebar'),
        'sidebar-soft': withVar('--sidebar-soft'),
        'sidebar-text': withVar('--sidebar-text'),
        cut: withVar('--cut'),
        'cut-dark': withVar('--cut-dark'),
        teal: withVar('--teal'),
        'st-ok-bg': withVar('--ok-bg'),
        'st-ok-fg': withVar('--ok-fg'),
        'st-warn-bg': withVar('--warn-bg'),
        'st-warn-fg': withVar('--warn-fg'),
        'st-info-bg': withVar('--info-bg'),
        'st-info-fg': withVar('--info-fg'),
        'st-prod-bg': withVar('--prod-bg'),
        'st-prod-fg': withVar('--prod-fg'),
        'st-done-bg': withVar('--done-bg'),
        'st-done-fg': withVar('--done-fg'),
        'st-danger-bg': withVar('--danger-bg'),
        'st-danger-fg': withVar('--danger-fg'),
      },
      borderRadius: {
        lg: '16px',
        md: '12px',
        sm: '8px',
      },
      boxShadow: {
        card: '0 1px 2px rgba(32,36,43,.04), 0 6px 20px rgba(32,36,43,.05)',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0', transform: 'translateY(8px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
      animation: {
        fadeIn: 'fadeIn 0.4s ease-out both',
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic':
          'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
      },
    },
  },
  darkMode: 'class',
  plugins: [require('@tailwindcss/forms'), heroui()],
};
export default config;
