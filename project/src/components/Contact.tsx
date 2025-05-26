import React from 'react';
import { Mail, MessageSquare, ArrowRight } from 'lucide-react';

const Contact = () => {
  return (
    <section id="contact" className="py-24 bg-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
            Ready to Build on Your Domain?
          </h2>
          <p className="text-lg text-gray-600">
            Now that your domain is set up, you can start building your decentralized web presence
          </p>
        </div>

        <div className="max-w-5xl mx-auto">
          <div className="bg-gradient-to-br from-slate-900 to-indigo-900 rounded-2xl overflow-hidden shadow-xl">
            <div className="grid grid-cols-1 md:grid-cols-2">
              <div className="p-8 md:p-12 flex items-center">
                <div>
                  <h3 className="text-2xl font-bold text-white mb-4">
                    Ready to Get Started?
                  </h3>
                  <p className="text-slate-300 mb-6">
                    Your domain is now connected and working properly. The next step is to replace this test page with your actual website or application.
                  </p>
                  
                  <div className="space-y-4">
                    <div className="flex items-start">
                      <div className="flex-shrink-0 p-2 bg-indigo-500/20 rounded-lg">
                        <MessageSquare className="h-5 w-5 text-indigo-300" />
                      </div>
                      <div className="ml-4">
                        <h4 className="text-white font-medium">Join Our Community</h4>
                        <p className="text-slate-300 text-sm">Get help from our active community of Web3 developers</p>
                      </div>
                    </div>
                    
                    <div className="flex items-start">
                      <div className="flex-shrink-0 p-2 bg-indigo-500/20 rounded-lg">
                        <Mail className="h-5 w-5 text-indigo-300" />
                      </div>
                      <div className="ml-4">
                        <h4 className="text-white font-medium">Documentation</h4>
                        <p className="text-slate-300 text-sm">Explore our guides to learn how to build on decentralized infrastructure</p>
                      </div>
                    </div>
                  </div>
                  
                  <div className="mt-8">
                    <button className="bg-white text-indigo-900 hover:bg-gray-100 px-6 py-3 rounded-lg text-base font-medium transition-all transform hover:scale-105 hover:shadow-lg flex items-center">
                      <span>View Documentation</span>
                      <ArrowRight className="ml-2 h-5 w-5" />
                    </button>
                  </div>
                </div>
              </div>
              
              <div className="p-8 md:p-12 bg-white/10 backdrop-blur-sm">
                <form className="space-y-4">
                  <div>
                    <label htmlFor="name" className="block text-sm font-medium text-white mb-1">
                      Name
                    </label>
                    <input
                      type="text"
                      id="name"
                      className="w-full px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-white placeholder-white/50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                      placeholder="Your name"
                    />
                  </div>
                  
                  <div>
                    <label htmlFor="email" className="block text-sm font-medium text-white mb-1">
                      Email
                    </label>
                    <input
                      type="email"
                      id="email"
                      className="w-full px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-white placeholder-white/50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                      placeholder="your@email.com"
                    />
                  </div>
                  
                  <div>
                    <label htmlFor="message" className="block text-sm font-medium text-white mb-1">
                      Message
                    </label>
                    <textarea
                      id="message"
                      rows={4}
                      className="w-full px-4 py-3 rounded-lg bg-white/5 border border-white/10 text-white placeholder-white/50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                      placeholder="Your message..."
                    ></textarea>
                  </div>
                  
                  <button 
                    type="submit" 
                    className="w-full bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-lg text-base font-medium transition-all transform hover:scale-105 hover:shadow-lg"
                  >
                    Send Message
                  </button>
                </form>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Contact;