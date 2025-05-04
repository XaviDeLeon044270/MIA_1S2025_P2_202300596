'use client';
import { Editor } from '@monaco-editor/react';
import { useState, useRef } from 'react';
import { executeCommands } from "@/services/api";

export default function Home() {
  const [command, setCommand] = useState<string>('');
  const [output, setOutput] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const editorRef = useRef<any>(null);

  const handleEditorDidMount = (editor: any) => {
    editorRef.current = editor;
  };

  const handleExecute = async () => {
    if (!command.trim()) {
      setOutput("Por favor ingrese un comando");
      return;
    }

    setIsLoading(true);
    try {
      const result = await executeCommands(command);
      setOutput(result);
    } catch (error) {
      setOutput(error instanceof Error ? error.message : "Error desconocido");
    } finally {
      setIsLoading(false);
    }
  };

  const handleOpenFile = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith('.smia')) {
      setOutput('Error: El archivo debe tener extensión .smia');
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      setCommand(content);
      
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }

      // Enfocar el editor después de cargar el archivo
      if (editorRef.current) {
        setTimeout(() => {
          editorRef.current.focus();
        }, 100);
      }
    };
    reader.onerror = () => {
      setOutput('Error al leer el archivo');
    };
    reader.readAsText(file);
  };

  const handleSaveFile = () => {
    if (!command.trim()) {
      setOutput("No hay contenido para guardar");
      return;
    }

    const blob = new Blob([command], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'archivo.smia';
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleNewFile = () => {
    if (command && !confirm('¿Deseas guardar el archivo antes de crear uno nuevo?')) {
      return;
    }
    setCommand('');
    setOutput('');

    // Enfocar el editor después de crear nuevo archivo
    if (editorRef.current) {
      setTimeout(() => {
        editorRef.current.focus();
      }, 100);
    }
  };

  return (
    <div className='flex flex-col min-h-screen bg-gray-900 text-white'>
      {/* Barra superior */}
      <div className='bg-gray-800 p-2 flex justify-between items-center'>
        <div className='flex space-x-4'>
          {/* Menú Archivo */}
          <div className='relative group'>
            <button className='bg-gray-700 hover:bg-gray-600 text-white font-bold py-2 px-6 rounded w-65'>
              Archivo
            </button>
            <div className='absolute hidden group-hover:block bg-gray-700 mt-1 rounded z-50 w-65'>
              <input
                type="file"
                id="file-input"
                className="hidden"
                accept=".smia"
                onChange={handleOpenFile}
                ref={fileInputRef}
              />
              <label
                htmlFor="file-input"
                className='block px-4 py-2 hover:bg-gray-600 cursor-pointer'
              >
                Abrir Archivo
              </label>
              <button
                onClick={handleSaveFile}
                className='block w-full px-4 py-2 hover:bg-gray-600 text-left'
              >
                Guardar Archivo
              </button>
              <button
                onClick={handleNewFile}
                className='block w-full px-4 py-2 hover:bg-gray-600 text-left'
              >
                Crear Archivo
              </button>
            </div>
          </div>
        </div>

        {/* Botón Run */}
        <button
          onClick={handleExecute}
          disabled={isLoading}
          className={`bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded
            ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
        >
          {isLoading ? (
            <div className="flex items-center">
              <svg className="animate-spin h-5 w-5 mr-2" viewBox="0 0 24 24">
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                  fill="none"
                />
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
              </svg>
              Ejecutando...
            </div>
          ) : (
            "Run"
          )}
        </button>
      </div>

      {/* Contenido principal */}
      <div className='flex flex-1 p-4 space-x-4'>
        {/* Editor */}
        <div className='flex-1 flex flex-col'>
          <label className='text-sm font-bold mb-2'>Editor:</label>
          <Editor
            height="70vh"
            defaultLanguage="go"
            theme='vs-dark'
            value={command}
            onChange={(value) => setCommand(value || '')}
            onMount={handleEditorDidMount}
            options={{
              minimap: { enabled: false },
              fontSize: 14,
              wordWrap: 'on',
              automaticLayout: true,
              lineNumbers: 'on',
              scrollBeyondLastLine: false,
              renderWhitespace: 'selection',
              tabSize: 2,
            }}
          />
        </div>

        {/* Consola */}
        <div className='flex-1 flex flex-col'>
          <label className='text-sm font-bold mb-2'>Consola:</label>
          <div className='bg-gray-800 p-4 rounded flex-1 overflow-auto'>
            <pre className='text-white'>{output}</pre>
          </div>
        </div>
      </div>
    </div>
  );
}