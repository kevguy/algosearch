import { Button } from "@mui/material";
import React, { useState } from "react";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import CheckIcon from "@mui/icons-material/Check";

const CopyIcon = ({
  copyableText,
  className,
}: {
  copyableText: string;
  className: string;
}) => {
  const [copied, setCopied] = useState(false);

  const copy = () => {
    if (navigator.clipboard && window.isSecureContext) {
      navigator.clipboard.writeText(copyableText);
      setCopied(true);
      setTimeout(() => {
        setCopied(false);
      }, 1000);
    }
  };

  return (
    <Button className={className} onClick={copy} size="small" aria-label="copy">
      {copied ? (
        <CheckIcon fontSize="small" />
      ) : (
        <ContentCopyIcon fontSize="small" />
      )}
    </Button>
  );
};

export default CopyIcon;
