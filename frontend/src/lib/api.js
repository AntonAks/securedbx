/**
 * Centralized API configuration
 * Replaces sed-based /prod ↔ /dev patching
 */

export const API_BASE = import.meta.env.VITE_API_BASE || '/prod';
export const RECAPTCHA_SITE_KEY = '6LdulTIsAAAAAJdhvyMU6B1og7GE7d5DySrQUQiv';
export const MAX_FILE_SIZE = 500 * 1024 * 1024; // 500 MB
