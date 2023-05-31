/*
Language: Odin
Author: Jon Lipstate <Jon@Lipstate.com>
Description: Odin-Language definition
Website: https://odin-lang.org
*/
export default function (hljs) {
	const ODIN_KEYWORDS = {
		keyword:
			'asm auto_cast bit_set break case cast context continue defer distinct do dynamic else enum ' +
			'fallthrough for foreign if import in map not_in or_else or_return package proc return struct ' +
			'switch transmute typeid union using when where ->',
		literal: 'true false nil ---',
		built_in:
			'int i8 i16 i32 i64 i128 ' +
			'uint byte u8 u16 u32 u64 u128 uintptr ' +
			'bool b8 b16 b32 b64 ' +
			'f16 f32 f64 ' +
			'i16le i32le i64le i128le u16le u32le u64le u128le ' +
			'i16be i32be i64be i128be u16be u32be u64be u128be ' +
			'f16le f32le f64le ' +
			'f16be f32be f64be ' +
			'complex32 complex64 complex128 ' +
			'quaternion64 quaternion128 quaternion256 ' +
			'rune string cstring any rawptr'
	};

	const ODIN_NUMBER_MODE = {
		className: 'number',
		begin:
			'\\b([0-9](_?[0-9])*|(0x|0o|0b)[0-9a-fA-F](_?[0-9a-fA-F])*)(\\.([0-9](_?[0-9])*)?)?([eE][-+]?[0-9](_?[0-9])*)?([pP][-+]?[0-9](_?[0-9])*)?[iufb]?',
		relevance: 0
	};
	const FUNCTION_DECL = {
		className: 'function',
		begin: /(\w+)\s*::\s*proc/,
		end: /$/,
		returnBegin: true,
		contains: [
			{
				className: 'title',
				begin: /(\w+)(?=\s*::)/
			}
		]
	};
	/* this is problematic as it shadows many other rules:*/
	const FUNCTION_CALL = {
		className: 'function',
		begin: /(\w+)\s*\(/,
		end: /$/,
		returnBegin: true,
		contains: [
			{
				className: 'title',
				begin: /(\w+)(?=\s*\()/
			},
			hljs.C_LINE_COMMENT_MODE,
			hljs.C_BLOCK_COMMENT_MODE,
			hljs.APOS_STRING_MODE,
			hljs.QUOTE_STRING_MODE,
			ODIN_NUMBER_MODE
		]
	};
	const RAW_STRING_MODE = {
		className: 'string',
		begin: '`', // Starts with `
		end: '`', // Ends with `
		contains: [],
		relevance: 0
	};

	return {
		name: 'Odin',
		aliases: ['odin', 'language-odin'],
		keywords: ODIN_KEYWORDS,
		contains: [
			hljs.C_LINE_COMMENT_MODE,
			hljs.C_BLOCK_COMMENT_MODE,
			hljs.APOS_STRING_MODE,
			hljs.QUOTE_STRING_MODE,
			RAW_STRING_MODE,
			ODIN_NUMBER_MODE
			// FUNCTION_CALL,
			// FUNCTION_DECL
		]
	};
}
