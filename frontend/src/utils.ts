/**
 * Formats bytes into a human-readable string with appropriate unit (B, KB, MB, GB, TB, PB)
 * @param bytes - The number of bytes to format
 * @param decimals - Number of decimal places (default: 2)
 * @returns Formatted string with unit
 */
export function formatBytes(bytes: number, decimals: number = 2): string {
	if (bytes === 0) return "0 B";

	const k = 1024;
	const dm = decimals < 0 ? 0 : decimals;
	const sizes = ["B", "KB", "MB", "GB", "TB", "PB"];

	const i = Math.floor(Math.log(bytes) / Math.log(k));

	return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + sizes[i];
}
