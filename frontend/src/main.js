import * as connectionPage from './pages/connection/script.js';
import * as nav from './pages/nav/script.js';
import * as graphPage from './pages/graph/script.js';
import * as integrityPage from './pages/integrity/script.js';
import './style.css';
import './app.css';

export const pagesKeys = {
    'connection': 'connection',
    'graph': 'graph',
    'integrity': 'integrity'
};

export const loadPage = async (pageName, data = undefined) => {
    try {
        const appDiv = document.getElementById('app');
        const navDiv = document.getElementById('nav');
        navDiv.style.display = 'none';

        // Pick the right module!
        let pageModule;
        switch (pageName) {
            case pagesKeys.connection:
                pageModule = connectionPage;
                break;
            case pagesKeys.graph:
                pageModule = graphPage;
                break;
            case pagesKeys.integrity:
                pageModule = integrityPage;
                break;
            default:
                pageModule = connectionPage;
                break;
        }

        if (pageModule) {
            const translations = await pageModule.getTranslations() ?? {};
            const pageElement = await injectTranslationsToHtmlString(pageModule.html, translations);
            appDiv.innerHTML = '';
            appDiv.appendChild(pageElement); // Insert the *initial* translated content into the DOM
            navDiv.style.display = pageName !== pagesKeys.connection ? 'flex' : 'none';

            /** 
             * Create and start the DOM observer for dynamic translations
             *  
             * https://hacks.mozilla.org/2012/05/dom-mutationobserver-reacting-to-dom-changes-without-killing-browser-performance/
            */ 
            const domObserver = createDomObserver(translations, appDiv);
            const param = data ?? {};
            param['pageName'] = pageName;
            param['observer'] = await domObserver;

            if (pageModule.init)  await pageModule.init(param);
            if (pageName !== pagesKeys.connection) nav.init(param);
        } else {
            console.error(`Page module for "${pageName}" not found.`);
        }
    } catch (error) {
        console.error(`Error loading page "${pageName}":`, error);
    }
}

// Landing page
document.addEventListener('DOMContentLoaded', async () => {
    await loadPage(pagesKeys.connection);
});

// Inject translated string vars to an HTML string and return as an element
export async function injectTranslationsToHtmlString(htmlString, translations) {
    // Create a tempo elem to contain HTML
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = htmlString;
    // Call trads on that tempo elem
    await injectTranslations(tempDiv, translations);

    return tempDiv.firstElementChild;
}

// Inject translated string vars to an HTML element
export async function injectTranslations(element, translations) {
    let missingTranslations = [];
    const unusedTranslations = new Set(Object.keys(translations)); // Using Set for faster search
    let translationsPerformed = false;

    // Text injection function
    const replaceText = (node) => {
        if (node.nodeType === Node.TEXT_NODE) {
            const originalValue = node.nodeValue;
            node.nodeValue = node.nodeValue.replace(/string:([a-zA-Z0-9_]+);/g, (match, key) => {
                if (translations.hasOwnProperty(key)) {
                    translationsPerformed = true;
                    unusedTranslations.delete(key);
                    return translations[key];
                } else {
                    missingTranslations.push(key);
                    return match;
                }
            });
            if (node.nodeValue !== originalValue) {
                translationsPerformed = true; // Marks trads has been used
            }
        } else if (node.nodeType === Node.ELEMENT_NODE) {
            node.childNodes.forEach(replaceText);
        }
    }

    // Applying text injection recursively
    replaceText(element);
    if (translationsPerformed) {
        /**
         * Logs text injection problems
         * this doesn't work fully since every module calls it 3x
         * but lets you see what was in all 3 times...
         */
        if (missingTranslations.length > 0) {
            console.log("Missing translations for keys:", missingTranslations);
        }
        if (unusedTranslations.size > 0) {
            console.log("Unused translations for keys:", Array.from(unusedTranslations));
        }
    }
}

async function createDomObserver(translations, appDiv) {
    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
                mutation.addedNodes.forEach(async (node) => {
                    if (node.nodeType === Node.ELEMENT_NODE) {
                        await injectTranslations(node, translations);
                    }
                });
            }
        });
    });

    const config = { childList: true, subtree: true };
    observer.observe(appDiv, config);

    return observer;
}

