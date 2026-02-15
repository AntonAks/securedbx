import { ref } from 'vue';

/**
 * Composable for file drag-and-drop handling
 */
export function useFileDrop(onFiles) {
    const isDragOver = ref(false);

    function handleDragOver(e) {
        e.preventDefault();
        isDragOver.value = true;
    }

    function handleDragLeave(e) {
        e.preventDefault();
        isDragOver.value = false;
    }

    function handleDrop(e) {
        e.preventDefault();
        isDragOver.value = false;
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            onFiles(files);
        }
    }

    return {
        isDragOver,
        handleDragOver,
        handleDragLeave,
        handleDrop,
    };
}
