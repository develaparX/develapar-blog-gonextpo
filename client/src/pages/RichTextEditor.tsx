import React from "react";
import { LexicalComposer } from "@lexical/react/LexicalComposer";
import { RichTextPlugin } from "@lexical/react/LexicalRichTextPlugin";
import { ContentEditable } from "@lexical/react/LexicalContentEditable";
import { HistoryPlugin } from "@lexical/react/LexicalHistoryPlugin";
import { AutoFocusPlugin } from "@lexical/react/LexicalAutoFocusPlugin";
import { LexicalErrorBoundary } from "@lexical/react/LexicalErrorBoundary";
import { HeadingNode, QuoteNode } from "@lexical/rich-text";
import { TableCellNode, TableNode, TableRowNode } from "@lexical/table";
import { ListItemNode, ListNode } from "@lexical/list";
import { CodeHighlightNode, CodeNode } from "@lexical/code";
import { AutoLinkNode, LinkNode } from "@lexical/link";
import { LinkPlugin } from "@lexical/react/LexicalLinkPlugin";
import { ListPlugin } from "@lexical/react/LexicalListPlugin";
import { OnChangePlugin } from "@lexical/react/LexicalOnChangePlugin";

interface RichTextEditorProps {
  onChange: (content: string) => void;
  initialContent?: string;
}

const RichTextEditor: React.FC<RichTextEditorProps> = ({
  onChange,
  initialContent,
}) => {
  // Theme untuk editor
  const editorTheme = {
    root: "bg-white rounded-md border border-gray-300 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 p-2 min-h-[300px]",
    paragraph: "my-2",
    heading: {
      h1: "text-2xl font-bold",
      h2: "text-xl font-bold",
      h3: "text-lg font-bold",
    },
    list: {
      ul: "list-disc ml-5",
      ol: "list-decimal ml-5",
    },
    link: "text-indigo-600 underline",
  };

  // Konfigurasi awal editor
  const initialConfig = {
    namespace: "ArticleEditor",
    theme: editorTheme,
    onError: (error: Error) => {
      console.error("Lexical Editor error:", error);
    },
    nodes: [
      HeadingNode,
      QuoteNode,
      ListItemNode,
      ListNode,
      TableNode,
      TableCellNode,
      TableRowNode,
      CodeNode,
      CodeHighlightNode,
      AutoLinkNode,
      LinkNode,
    ],
  };

  // Handler untuk perubahan konten
  const handleEditorChange = (editorState: any) => {
    editorState.read(() => {
      const jsonString = JSON.stringify(editorState.toJSON());
      onChange(jsonString);
    });
  };

  return (
    <LexicalComposer initialConfig={initialConfig}>
      <div className="editor-container">
        <RichTextPlugin
          contentEditable={<ContentEditable className="editor-input" />}
          placeholder={
            <div className="text-gray-400 absolute top-3 left-3">
              Tulis konten artikel Anda di sini...
            </div>
          }
          ErrorBoundary={LexicalErrorBoundary}
        />
        <HistoryPlugin />
        <AutoFocusPlugin />
        <ListPlugin />
        <LinkPlugin />
        <OnChangePlugin onChange={handleEditorChange} />
      </div>
    </LexicalComposer>
  );
};

export default RichTextEditor;
