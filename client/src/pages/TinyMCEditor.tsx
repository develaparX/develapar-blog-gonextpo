import React, { useRef } from "react";
import { Editor } from "@tinymce/tinymce-react";

interface TinyMCEEditorProps {
  value: string;
  onChange: (content: string) => void;
  height?: number;
}

const TinyMCEEditor: React.FC<TinyMCEEditorProps> = ({
  value,
  onChange,
  height = 400,
}) => {
  const editorRef = useRef<any>(null);

  const handleEditorChange = (content: string) => {
    onChange(content);
  };

  return (
    <Editor
      apiKey="8i9m8wgxfdo8yp6jgi81kngnf44sjpujhj384vr85mn2a3ev" // Gunakan API key gratis dari TinyMCE Cloud
      onInit={(evt, editor): any => (editorRef.current = editor)}
      value={value}
      onEditorChange={handleEditorChange}
      init={{
        height: height,
        menubar: true,
        plugins: [
          "advlist",
          "autolink",
          "lists",
          "link",
          "image",
          "charmap",
          "preview",
          "anchor",
          "searchreplace",
          "visualblocks",
          "code",
          "fullscreen",
          "insertdatetime",
          "media",
          "table",
          "code",
          "help",
          "wordcount",
          "emoticons",
          "template",
          "paste",
          "textcolor",
          "colorpicker",
        ],
        toolbar:
          "undo redo | blocks fontfamily fontsize | " +
          "bold italic underline strikethrough | link image media table | " +
          "alignleft aligncenter alignright alignjustify | " +
          "numlist bullist outdent indent | forecolor backcolor | " +
          "emoticons charmap | insertdatetime | preview code | help",
        content_style: `
          body { 
            font-family: Helvetica,Arial,sans-serif; 
            font-size: 14px;
            line-height: 1.6;
          }
        `,
        paste_data_images: true,
        image_advtab: true,
        templates: [
          {
            title: "New Table",
            description: "creates a new table",
            content:
              '<div class="mceTmpl"><table width="98%%"  border="0" cellspacing="0" cellpadding="0"><tr><th scope="col"> </th><th scope="col"> </th></tr><tr><td> </td><td> </td></tr></table></div>',
          },
          {
            title: "Starting my story",
            description: "A cure for writers block",
            content: "Once upon a time...",
          },
        ],
        template_cdate_format: "[Date Created (CDATE): %m/%d/%Y : %H:%M:%S]",
        template_mdate_format: "[Date Modified (MDATE): %m/%d/%Y : %H:%M:%S]",
        image_caption: true,
        quickbars_selection_toolbar:
          "bold italic | quicklink h2 h3 blockquote quickimage quicktable",
        noneditable_noneditable_class: "mceNonEditable",
        toolbar_mode: "sliding",
        contextmenu: "link image imagetools table",
      }}
    />
  );
};

export default TinyMCEEditor;
