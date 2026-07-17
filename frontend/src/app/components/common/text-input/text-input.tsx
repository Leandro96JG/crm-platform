'use client';

import { HTMLInputTypeAttribute } from "react";

interface TextInputProps {
  name: string;
  label: string;
  type?: HTMLInputTypeAttribute;
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  className?: string;
  defaultValue?: string;
}

export function TextInput({ label, name, type, placeholder, required, className, defaultValue, disabled }: TextInputProps) {
  return (
    <div className={`${className ? className : ''}`}>
      <label
        className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-ink-faint"
        htmlFor={name}
      >
        {label}
      </label>

      <input
        className="peer block w-full rounded-md border border-line bg-card px-3 py-2.5 text-sm text-ink outline-none transition-colors placeholder:text-ink-faint focus:border-cut focus:ring-1 focus:ring-cut disabled:cursor-not-allowed disabled:opacity-60"
        id={name}
        type={type || "text"}
        name={name}
        placeholder={placeholder}
        required={required}
        defaultValue={defaultValue}
        disabled={disabled}
      />
    </div>
  );
}