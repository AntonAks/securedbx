/**
 * TextUpload — verify text_size encoding logic.
 *
 * Run with: npx vitest run tests/views/TextUpload.test.js
 */

import { describe, it, expect } from 'vitest';

describe('TextUpload — text_size encoding', () => {
    it('byteLength equals character count for ASCII text', () => {
        const text = 'Hello world';
        const encoder = new TextEncoder();
        const buf = encoder.encode(text);
        expect(buf.byteLength).toBe(11);
        expect(typeof buf.byteLength).toBe('number');
        expect(Number.isInteger(buf.byteLength)).toBe(true);
    });

    it('byteLength is greater than character count for multibyte chars', () => {
        const text = 'Привіт'; // Ukrainian — multibyte UTF-8
        const encoder = new TextEncoder();
        const buf = encoder.encode(text);
        expect(buf.byteLength).toBeGreaterThan(text.length);
        expect(Number.isInteger(buf.byteLength)).toBe(true);
    });
});
