import React from "react";
import styles from "./HomeFooter.module.scss";

const HomeFooter = () => {
  return (
    <div className={styles.homefooter}>
      <div>
        <h3>About Algorand</h3>
        <ul>
          <li>
            <a
              href="https://www.algorand.com/"
              target="_blank"
              rel="noopener noreferrer"
            >
              Algorand Inc
            </a>
          </li>
          <li>
            <a
              href="https://algorand.foundation/"
              target="_blank"
              rel="noopener noreferrer"
            >
              Algorand Foundation
            </a>
          </li>
          <li>
            <a
              href="https://www.algorand.com/technology/faq"
              target="_blank"
              rel="noopener noreferrer"
            >
              Frequently Asked Questions
            </a>
          </li>
        </ul>
      </div>
      <div>
        <h3>Getting Started</h3>
        <ul>
          <li>
            <a
              href="https://developer.algorand.org/"
              target="_blank"
              rel="noopener noreferrer"
            >
              Documentation and Tutorials
            </a>
          </li>
          <li>
            <a
              href="https://github.com/algorand"
              target="_blank"
              rel="noopener noreferrer"
            >
              Contribute to Open Source
            </a>
          </li>
          <li>
            <a
              href="https://forum.algorand.org/"
              target="_blank"
              rel="noopener noreferrer"
            >
              Join the Discussion
            </a>
          </li>
        </ul>
      </div>
      <div>
        <h3>Technology</h3>
        <ul>
          <li>
            <a
              href="https://www.algorand.com/technology"
              target="_blank"
              rel="noopener noreferrer"
            >
              Algorand Features &amp; Capabilities in Layer-1
            </a>
          </li>
          <li>
            <a
              href="https://www.algorand.com/technology/research-innovation"
              target="_blank"
              rel="noopener noreferrer"
            >
              Research and Innovation
            </a>
          </li>
          <li>
            <a
              href="https://www.algorand.com/resources/white-papers"
              target="_blank"
              rel="noopener noreferrer"
            >
              Official White Papers
            </a>
          </li>
        </ul>
      </div>
    </div>
  );
};

export default HomeFooter;
