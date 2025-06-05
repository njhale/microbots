from fastmcp import FastMCP
from markitdown import MarkItDown
import requests
from typing import Optional, List
from bs4 import BeautifulSoup

mcp = FastMCP("hackerdigest")
md = MarkItDown()


@mcp.tool()
async def summarize_link(link: str) -> str:
    """Fetches a link and returns a summary of its content."""
    try:
        response = requests.get(link, timeout=10)
        response.raise_for_status()
        
        # Convert to markdown and return first 500 characters as summary
        markdown_content = md.convert(response).text_content
        summary = markdown_content[:500] + "..." if len(markdown_content) > 500 else markdown_content
        return f"Summary of {link}:\n{summary}"
    except Exception as e:
        return f"Error fetching {link}: {str(e)}"

@mcp.tool()
async def get_hackernews_links() -> List[str]:
    """Gets the top 5 links from Hacker News home page."""
    try:
        response = requests.get("https://news.ycombinator.com/", timeout=10)
        response.raise_for_status()
        
        soup = BeautifulSoup(response.content, 'html.parser')
        
        # Find all story links (they have class 'storylink' or are in 'titleline' spans)
        links = []
        storylinks = soup.find_all('span', class_='titleline')
        
        for storylink in storylinks[:5]:  # Get top 5
            link_element = storylink.find('a')
            if link_element and link_element.get('href'):
                href = link_element.get('href')
                # Handle relative URLs
                if href.startswith('item?'):
                    href = f"https://news.ycombinator.com/{href}"
                elif not href.startswith('http'):
                    href = f"https://news.ycombinator.com/{href}"
                links.append(href)
        
        return links
    except Exception as e:
        return [f"Error fetching Hacker News: {str(e)}"]

@mcp.tool()
async def fetch_and_convert_links(links: List[str]) -> str:
    """Fetches a list of links and converts their content to markdown."""
    results = []
    
    for i, link in enumerate(links, 1):
        try:
            response = requests.get(link, timeout=10)
            response.raise_for_status()
            
            # Convert to markdown
            markdown_content = md.convert(response).text_content
            
            results.append(f"## Link {i}: {link}\n\n{markdown_content}\n\n---\n")
            
        except Exception as e:
            results.append(f"## Link {i}: {link}\n\nError: {str(e)}\n\n---\n")
    
    return "\n".join(results)


if __name__ == "__main__":
    mcp.run()